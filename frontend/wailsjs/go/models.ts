export namespace main {
	
	export class DiskTree {
	    name: string;
	    children: any[];
	
	    static createFrom(source: any = {}) {
	        return new DiskTree(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.children = source["children"];
	    }
	}

}

