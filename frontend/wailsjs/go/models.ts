export namespace main {
	
	export class DeviceInfo {
	    device_id: string;
	    device_model: string;
	    device_name: string;
	    os_version: string;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.device_id = source["device_id"];
	        this.device_model = source["device_model"];
	        this.device_name = source["device_name"];
	        this.os_version = source["os_version"];
	    }
	}

}

