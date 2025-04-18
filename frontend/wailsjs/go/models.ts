export namespace main {
	
	export class LauncherSettings {
	    gameDirectory: string;
	    allocatedRAM?: number;
	    jvmArguments?: string;
	    showAlpha: boolean;
	    showBeta: boolean;
	    showSnapshots: boolean;
	    showOldVersions: boolean;
	    showOnlyInstalled: boolean;
	    resolutionWidth?: number;
	    resolutionHeight?: number;
	
	    static createFrom(source: any = {}) {
	        return new LauncherSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameDirectory = source["gameDirectory"];
	        this.allocatedRAM = source["allocatedRAM"];
	        this.jvmArguments = source["jvmArguments"];
	        this.showAlpha = source["showAlpha"];
	        this.showBeta = source["showBeta"];
	        this.showSnapshots = source["showSnapshots"];
	        this.showOldVersions = source["showOldVersions"];
	        this.showOnlyInstalled = source["showOnlyInstalled"];
	        this.resolutionWidth = source["resolutionWidth"];
	        this.resolutionHeight = source["resolutionHeight"];
	    }
	}
	export class StartOptions {
	    version?: minecraft.MinecraftVersionInfo;
	
	    static createFrom(source: any = {}) {
	        return new StartOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = this.convertValues(source["version"], minecraft.MinecraftVersionInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace minecraft {
	
	export class MinecraftVersionInfo {
	    id: string;
	    type: string;
	    releaseTime: string;
	    complianceLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new MinecraftVersionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.releaseTime = source["releaseTime"];
	        this.complianceLevel = source["complianceLevel"];
	    }
	}

}

