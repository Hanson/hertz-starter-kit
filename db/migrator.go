package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
	"regexp"
	"strings"
)

// 目的：string 默认转为 varchar(100), 所有字段 not null, int 默认值为 0
// 搜索 @修改点 查找修改源码的地方

type MyMigrator struct {
	mysql.Migrator
}

var regFullDataType = regexp.MustCompile(`\D*(\d+)\D?`)

type printSQLLogger struct {
	logger.Interface
}

func CustomizeField(field *schema.Field) {
	if strings.Contains(strings.ToLower(string(field.DataType)), "text") {
		field.NotNull = false
	} else {
		field.NotNull = true
	}

	if field.Name == "DeletedAt" {
		field.HasDefaultValue = true
		field.DefaultValue = "0"
	} else if field.DataType == "string" || strings.HasPrefix(strings.ToLower(string(field.DataType)), "varchar") {
		field.HasDefaultValue = true
		field.DefaultValue = ""
		field.DefaultValueInterface = ""
	}

	if field.DataType == "string" {
		field.DataType = "varchar(100)"
	}
}

// AutoMigrate auto migrate values
func (m MyMigrator) AutoMigrate(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, true) {
		queryTx := m.DB.Session(&gorm.Session{})
		execTx := queryTx
		if m.DB.DryRun {
			queryTx.DryRun = false
			execTx = m.DB.Session(&gorm.Session{Logger: &printSQLLogger{Interface: m.DB.Logger}})
		}
		if !queryTx.Migrator().HasTable(value) {
			// @修改点
			if err := m.CreateTable(value); err != nil {
				return err
			}
		} else {
			if err := m.RunWithValue(value, func(stmt *gorm.Statement) (errr error) {
				columnTypes, err := queryTx.Migrator().ColumnTypes(value)
				if err != nil {
					return err
				}
				var (
					parseIndexes          = stmt.Schema.ParseIndexes()
					parseCheckConstraints = stmt.Schema.ParseCheckConstraints()
				)
				for _, dbName := range stmt.Schema.DBNames {
					field := stmt.Schema.FieldsByDBName[dbName]
					var foundColumn gorm.ColumnType

					for _, columnType := range columnTypes {
						if columnType.Name() == dbName {
							foundColumn = columnType
							break
						}
					}
					if foundColumn == nil {
						// not found, add column
						if err := m.AddColumn(value, dbName); err != nil {
							return err
						}
						// @修改点
					} else if err := m.MigrateColumn(value, field, foundColumn); err != nil {
						// found, smart migrate
						return err
					}
				}

				if !m.DB.DisableForeignKeyConstraintWhenMigrating && !m.DB.IgnoreRelationshipsWhenMigrating {
					for _, rel := range stmt.Schema.Relationships.Relations {
						if rel.Field.IgnoreMigration {
							continue
						}
						if constraint := rel.ParseConstraint(); constraint != nil &&
							constraint.Schema == stmt.Schema && !queryTx.Migrator().HasConstraint(value, constraint.Name) {
							if err := execTx.Migrator().CreateConstraint(value, constraint.Name); err != nil {
								return err
							}
						}
					}
				}

				for _, chk := range parseCheckConstraints {
					if !queryTx.Migrator().HasConstraint(value, chk.Name) {
						if err := execTx.Migrator().CreateConstraint(value, chk.Name); err != nil {
							return err
						}
					}
				}

				for _, idx := range parseIndexes {
					if !queryTx.Migrator().HasIndex(value, idx.Name) {
						if err := execTx.Migrator().CreateIndex(value, idx.Name); err != nil {
							return err
						}
					}
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m MyMigrator) MigrateColumn(value interface{}, field *schema.Field, columnType gorm.ColumnType) error {
	// @修改点
	CustomizeField(field)

	// found, smart migrate
	fullDataType := strings.TrimSpace(strings.ToLower(m.DB.Migrator().FullDataTypeOf(field).SQL))
	realDataType := strings.ToLower(columnType.DatabaseTypeName())

	var (
		alterColumn bool
		isSameType  = fullDataType == realDataType
	)

	if !field.PrimaryKey {
		// check type
		if !strings.HasPrefix(fullDataType, realDataType) {
			// check type aliases
			aliases := m.DB.Migrator().GetTypeAliases(realDataType)
			for _, alias := range aliases {
				if strings.HasPrefix(fullDataType, alias) {
					isSameType = true
					break
				}
			}

			if !isSameType {
				alterColumn = true

			}
		}
	}

	if !isSameType {
		// check size
		if length, ok := columnType.Length(); length != int64(field.Size) {
			if length > 0 && field.Size > 0 {
				alterColumn = true

			} else {
				// has size in data type and not equal
				// Since the following code is frequently called in the for loop, reg optimization is needed here
				matches2 := regFullDataType.FindAllStringSubmatch(fullDataType, -1)
				if !field.PrimaryKey &&
					(len(matches2) == 1 && matches2[0][1] != fmt.Sprint(length) && ok) {
					alterColumn = true

				}
			}
		}

		// check precision
		if precision, _, ok := columnType.DecimalSize(); ok && int64(field.Precision) != precision {
			if regexp.MustCompile(fmt.Sprintf("[^0-9]%d[^0-9]", field.Precision)).MatchString(m.Migrator.Migrator.DataTypeOf(field)) {
				alterColumn = true

			}
		}
	}

	// check nullable
	if nullable, ok := columnType.Nullable(); ok && nullable == field.NotNull {
		// not primary key & database is nullable
		if !field.PrimaryKey && nullable {
			alterColumn = true

		}
	}

	// check unique
	if unique, ok := columnType.Unique(); ok && unique != field.Unique {
		// not primary key
		if !field.PrimaryKey {
			alterColumn = true

		}
	}

	// check default value
	if !field.PrimaryKey {
		currentDefaultNotNull := field.HasDefaultValue && (field.DefaultValueInterface != nil || !strings.EqualFold(field.DefaultValue, "NULL"))
		dv, dvNotNull := columnType.DefaultValue()
		if dvNotNull && !currentDefaultNotNull {
			// defalut value -> null
			alterColumn = true

		} else if !dvNotNull && currentDefaultNotNull {
			// null -> default value
			alterColumn = true

		} else if (field.GORMDataType != schema.Time && dv != field.DefaultValue) ||
			(field.GORMDataType == schema.Time && !strings.EqualFold(strings.TrimSuffix(dv, "()"), strings.TrimSuffix(field.DefaultValue, "()"))) {
			// default value not equal
			// not both null
			if currentDefaultNotNull || dvNotNull {
				alterColumn = true

			}
		}
	}

	// check comment
	if comment, ok := columnType.Comment(); ok && comment != field.Comment {
		// not primary key
		if !field.PrimaryKey {
			alterColumn = true

		}
	}

	if alterColumn && !field.IgnoreMigration {
		return m.DB.Migrator().AlterColumn(value, field.DBName)
	}

	return nil
}

// CreateTable create table in database for values
func (m MyMigrator) CreateTable(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, false) {
		tx := m.DB.Session(&gorm.Session{})
		if err := m.RunWithValue(value, func(stmt *gorm.Statement) (errr error) {
			var (
				createTableSQL          = "CREATE TABLE ? ("
				values                  = []interface{}{m.CurrentTable(stmt)}
				hasPrimaryKeyInDataType bool
			)

			for _, dbName := range stmt.Schema.DBNames {
				field := stmt.Schema.FieldsByDBName[dbName]
				CustomizeField(field)
				if !field.IgnoreMigration {
					createTableSQL += "? ?"
					hasPrimaryKeyInDataType = hasPrimaryKeyInDataType || strings.Contains(strings.ToUpper(string(field.DataType)), "PRIMARY KEY")
					values = append(values, clause.Column{Name: dbName}, m.DB.Migrator().FullDataTypeOf(field))
					createTableSQL += ","
				}
			}

			if !hasPrimaryKeyInDataType && len(stmt.Schema.PrimaryFields) > 0 {
				createTableSQL += "PRIMARY KEY ?,"
				primaryKeys := []interface{}{}
				for _, field := range stmt.Schema.PrimaryFields {
					primaryKeys = append(primaryKeys, clause.Column{Name: field.DBName})
				}

				values = append(values, primaryKeys)
			}

			for _, idx := range stmt.Schema.ParseIndexes() {
				if m.CreateIndexAfterCreateTable {
					defer func(value interface{}, name string) {
						if errr == nil {
							errr = tx.Migrator().CreateIndex(value, name)
						}
					}(value, idx.Name)
				} else {
					if idx.Class != "" {
						createTableSQL += idx.Class + " "
					}
					createTableSQL += "INDEX ? ?"

					if idx.Comment != "" {
						createTableSQL += fmt.Sprintf(" COMMENT '%s'", idx.Comment)
					}

					if idx.Option != "" {
						createTableSQL += " " + idx.Option
					}

					createTableSQL += ","
					values = append(values, clause.Column{Name: idx.Name}, tx.Migrator().(migrator.BuildIndexOptionsInterface).BuildIndexOptions(idx.Fields, stmt))
				}
			}

			if !m.DB.DisableForeignKeyConstraintWhenMigrating && !m.DB.IgnoreRelationshipsWhenMigrating {
				for _, rel := range stmt.Schema.Relationships.Relations {
					if rel.Field.IgnoreMigration {
						continue
					}
					if constraint := rel.ParseConstraint(); constraint != nil {
						if constraint.Schema == stmt.Schema {
							sql, vars := buildConstraint(constraint)
							createTableSQL += sql + ","
							values = append(values, vars...)
						}
					}
				}
			}

			for _, chk := range stmt.Schema.ParseCheckConstraints() {
				createTableSQL += "CONSTRAINT ? CHECK (?),"
				values = append(values, clause.Column{Name: chk.Name}, clause.Expr{SQL: chk.Constraint})
			}

			createTableSQL = strings.TrimSuffix(createTableSQL, ",")

			createTableSQL += ")"

			if tableOption, ok := m.DB.Get("gorm:table_options"); ok {
				createTableSQL += fmt.Sprint(tableOption)
			}

			errr = tx.Exec(createTableSQL, values...).Error
			return errr
		}); err != nil {
			return err
		}
	}
	return nil
}

func buildConstraint(constraint *schema.Constraint) (sql string, results []interface{}) {
	sql = "CONSTRAINT ? FOREIGN KEY ? REFERENCES ??"
	if constraint.OnDelete != "" {
		sql += " ON DELETE " + constraint.OnDelete
	}

	if constraint.OnUpdate != "" {
		sql += " ON UPDATE " + constraint.OnUpdate
	}

	var foreignKeys, references []interface{}
	for _, field := range constraint.ForeignKeys {
		foreignKeys = append(foreignKeys, clause.Column{Name: field.DBName})
	}

	for _, field := range constraint.References {
		references = append(references, clause.Column{Name: field.DBName})
	}
	results = append(results, clause.Table{Name: constraint.Name}, foreignKeys, clause.Table{Name: constraint.ReferenceSchema.Table}, references)
	return
}

// AddColumn create `name` column for value
func (m MyMigrator) AddColumn(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		// avoid using the same name field
		f := stmt.Schema.LookUpField(name)
		if f == nil {
			return fmt.Errorf("failed to look up field with name: %s", name)
		}

		CustomizeField(f)

		if !f.IgnoreMigration {
			return m.DB.Exec(
				"ALTER TABLE ? ADD ? ?",
				m.CurrentTable(stmt), clause.Column{Name: f.DBName}, m.DB.Migrator().FullDataTypeOf(f),
			).Error
		}

		return nil
	})
}
