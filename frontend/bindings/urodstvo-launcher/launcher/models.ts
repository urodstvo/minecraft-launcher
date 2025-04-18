// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as minecraft$0 from "../minecraft/models.js";

export class AccountsInfo {
    "selectedAccount"?: string;
    "accounts"?: LauncherAccount[];

    /** Creates a new AccountsInfo instance. */
    constructor($$source: Partial<AccountsInfo> = {}) {

        Object.assign(this, $$source);
    }

    /**
     * Creates a new AccountsInfo instance from a string or object.
     */
    static createFrom($$source: any = {}): AccountsInfo {
        const $$createField1_0 = $$createType1;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("accounts" in $$parsedSource) {
            $$parsedSource["accounts"] = $$createField1_0($$parsedSource["accounts"]);
        }
        return new AccountsInfo($$parsedSource as Partial<AccountsInfo>);
    }
}

export class LauncherAccount {
    "id": string;
    "name": string;
    "skins": minecraft$0.MinecraftProfileSkin[];
    "capes": minecraft$0.MinecraftProfileCape[];
    "error": string;
    "errorMessage": string;
    "access_token": string;
    "refresh_token": string;

    /** Creates a new LauncherAccount instance. */
    constructor($$source: Partial<LauncherAccount> = {}) {
        if (!("id" in $$source)) {
            this["id"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("skins" in $$source)) {
            this["skins"] = [];
        }
        if (!("capes" in $$source)) {
            this["capes"] = [];
        }
        if (!("error" in $$source)) {
            this["error"] = "";
        }
        if (!("errorMessage" in $$source)) {
            this["errorMessage"] = "";
        }
        if (!("access_token" in $$source)) {
            this["access_token"] = "";
        }
        if (!("refresh_token" in $$source)) {
            this["refresh_token"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new LauncherAccount instance from a string or object.
     */
    static createFrom($$source: any = {}): LauncherAccount {
        const $$createField2_0 = $$createType3;
        const $$createField3_0 = $$createType5;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("skins" in $$parsedSource) {
            $$parsedSource["skins"] = $$createField2_0($$parsedSource["skins"]);
        }
        if ("capes" in $$parsedSource) {
            $$parsedSource["capes"] = $$createField3_0($$parsedSource["capes"]);
        }
        return new LauncherAccount($$parsedSource as Partial<LauncherAccount>);
    }
}

export class LauncherSettings {
    "gameDirectory": string;
    "allocatedRAM"?: number;
    "jvmArguments"?: string;
    "showAlpha": boolean;
    "showBeta": boolean;
    "showSnapshots": boolean;
    "showOldVersions": boolean;
    "showOnlyInstalled": boolean;
    "resolutionWidth"?: number;
    "resolutionHeight"?: number;

    /** Creates a new LauncherSettings instance. */
    constructor($$source: Partial<LauncherSettings> = {}) {
        if (!("gameDirectory" in $$source)) {
            this["gameDirectory"] = "";
        }
        if (!("showAlpha" in $$source)) {
            this["showAlpha"] = false;
        }
        if (!("showBeta" in $$source)) {
            this["showBeta"] = false;
        }
        if (!("showSnapshots" in $$source)) {
            this["showSnapshots"] = false;
        }
        if (!("showOldVersions" in $$source)) {
            this["showOldVersions"] = false;
        }
        if (!("showOnlyInstalled" in $$source)) {
            this["showOnlyInstalled"] = false;
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new LauncherSettings instance from a string or object.
     */
    static createFrom($$source: any = {}): LauncherSettings {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new LauncherSettings($$parsedSource as Partial<LauncherSettings>);
    }
}

export class StartOptions {
    "version": minecraft$0.MinecraftVersionInfo | null;

    /** Creates a new StartOptions instance. */
    constructor($$source: Partial<StartOptions> = {}) {
        if (!("version" in $$source)) {
            this["version"] = null;
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new StartOptions instance from a string or object.
     */
    static createFrom($$source: any = {}): StartOptions {
        const $$createField0_0 = $$createType7;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("version" in $$parsedSource) {
            $$parsedSource["version"] = $$createField0_0($$parsedSource["version"]);
        }
        return new StartOptions($$parsedSource as Partial<StartOptions>);
    }
}

// Private type creation functions
const $$createType0 = LauncherAccount.createFrom;
const $$createType1 = $Create.Array($$createType0);
const $$createType2 = minecraft$0.MinecraftProfileSkin.createFrom;
const $$createType3 = $Create.Array($$createType2);
const $$createType4 = minecraft$0.MinecraftProfileCape.createFrom;
const $$createType5 = $Create.Array($$createType4);
const $$createType6 = minecraft$0.MinecraftVersionInfo.createFrom;
const $$createType7 = $Create.Nullable($$createType6);
