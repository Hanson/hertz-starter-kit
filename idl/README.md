message CreateDemoReq {
	ModelDemo demo = 1[(api.vd)="$!=nil"];
}

message CreateDemoResp {
	ModelDemo demo = 1;
}

message UpdateDemoReq {
	uint64 id = 1[(api.vd)="$!=0"];

	map<string, string> update_map = 2[(api.vd)="$!=nil"];
}

message UpdateDemoResp {
}

message DeleteDemoReq {
	repeated uint64 id_list = 1;
}

message DeleteDemoResp {
}

message ListDemoReq {
	enum ListOption {
		ListOptionNil = 0;

		ListOptionDescription = 1;
	}
	paginate.ListOption list_option = 1[(api.vd)="$!=nil"];
}

message ListDemoResp {
	paginate.Paginate paginate = 1;

	repeated ModelDemo list = 2;
}

service DemoService {
	rpc CreateDemo(CreateDemoReq) returns(CreateDemoResp) {
		option (api.post) = "/admin/CreateDemo";
	}

	rpc UpdateDemo(UpdateDemoReq) returns(UpdateDemoResp) {
		option (api.post) = "/admin/UpdateDemo";
	}

	rpc DeleteDemo(DeleteDemoReq) returns(DeleteDemoResp) {
		option (api.post) = "/admin/DeleteDemo";
	}

	rpc ListDemo(ListDemoReq) returns(ListDemoResp) {
		option (api.post) = "/admin/ListDemo";
	}
}