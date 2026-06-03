import { CMCL_CONFIG } from '../variables';

const { shellApiUrl, cmclPath } = CMCL_CONFIG;

export interface ShellResult {
  success: boolean;
  stdout: string;
  stderr: string;
  error: string | null;
}

export interface InstallOptions {
  fabric?: boolean;
  fabricVersion?: string;
  forge?: boolean;
  forgeVersion?: string;
  optifine?: boolean;
  optifineVersion?: string;
  quilt?: boolean;
  quiltVersion?: string;
  select?: boolean;
}

export const ShellAPI = {
  async execute(command: string): Promise<ShellResult> {
    if (!command || typeof command !== 'string') {
      return { success: false, stdout: '', stderr: '', error: 'Invalid command' };
    }

    try {
      const response = await fetch(`${shellApiUrl}/shell`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ shell: command, type: 'simple' }),
      });

      if (!response.ok) {
        try {
          const json = await response.json();
          return {
            success: false,
            stdout: '',
            stderr: '',
            error: json.error || `HTTP ${response.status}: ${response.statusText}`,
          };
        } catch {
          const text = await response.text();
          return {
            success: false,
            stdout: '',
            stderr: '',
            error: `HTTP ${response.status}: ${text || response.statusText}`,
          };
        }
      }

      const result = await response.json();
      return {
        success: result.success || false,
        stdout: result.output || '',
        stderr: '',
        error: result.error || null,
      };
    } catch (err) {
      return {
        success: false,
        stdout: '',
        stderr: '',
        error: err instanceof Error ? err.message : 'Failed to execute command',
      };
    }
  },

  async ping(): Promise<{ success: boolean; responseTime: number; error?: string }> {
    const startTime = Date.now();
    try {
      const response = await fetch(`${shellApiUrl}/ping`, {
        method: 'GET',
      });
      const responseTime = Date.now() - startTime;
      if (response.ok) {
        return { success: true, responseTime };
      }
      return { success: false, responseTime, error: `HTTP ${response.status}: ${response.statusText}` };
    } catch (err) {
      const responseTime = Date.now() - startTime;
      return { success: false, responseTime, error: err instanceof Error ? err.message : 'Failed to ping' };
    }
  },
};

export const CMCLAPI = {
  async listVersions(): Promise<string[]> {
    const command = `${cmclPath} --list`;
    const result = await ShellAPI.execute(command);
    if (result.success) {
      const lines = result.stdout.trim().split('\n');
      return lines.filter(line => line.trim()).map(line => line.trim());
    }
    return [];
  },

  async selectVersion(version: string): Promise<ShellResult> {
    const command = `${cmclPath} --select=${version}`;
    return await ShellAPI.execute(command);
  },

  async startGame(version?: string): Promise<ShellResult> {
    const command = version ? `${cmclPath} ${version}` : cmclPath;
    return await ShellAPI.execute(command);
  },

  async installVersion(version: string, options: InstallOptions = {}): Promise<ShellResult> {
    let command = `${cmclPath} install ${version}`;
    
    if (options.fabric) {
      command += ` --fabric${options.fabricVersion ? `=${options.fabricVersion}` : ''}`;
    }
    if (options.forge) {
      command += ` --forge${options.forgeVersion ? `=${options.forgeVersion}` : ''}`;
    }
    if (options.optifine) {
      command += ` --optifine${options.optifineVersion ? `=${options.optifineVersion}` : ''}`;
    }
    if (options.quilt) {
      command += ` --quilt${options.quiltVersion ? `=${options.quiltVersion}` : ''}`;
    }
    if (options.select) {
      command += ' -s';
    }

    return await ShellAPI.execute(command);
  },

  async uninstallVersion(version: string): Promise<ShellResult> {
    const command = `${cmclPath} version ${version} --delete`;
    return await ShellAPI.execute(command);
  },

  async getVersionInfo(version: string): Promise<ShellResult> {
    const command = `${cmclPath} version ${version} --info`;
    return await ShellAPI.execute(command);
  },

  async listAccounts(): Promise<string[]> {
    const command = `${cmclPath} account --list`;
    const result = await ShellAPI.execute(command);
    if (result.success) {
      const lines = result.stdout.trim().split('\n');
      return lines.filter(line => line.trim()).map(line => line.trim());
    }
    return [];
  },

  async selectAccount(index: number): Promise<ShellResult> {
    const command = `${cmclPath} account --select=${index}`;
    return await ShellAPI.execute(command);
  },

  async loginOffline(username: string): Promise<ShellResult> {
    const command = `${cmclPath} account --login=offline --name=${username}`;
    return await ShellAPI.execute(command);
  },

  async loginMicrosoft(): Promise<ShellResult> {
    const command = `${cmclPath} account --login=microsoft`;
    return await ShellAPI.execute(command);
  },

  async getConfig(key?: string): Promise<ShellResult> {
    const command = key ? `${cmclPath} config ${key}` : `${cmclPath} config -a`;
    return await ShellAPI.execute(command);
  },

  async setConfig(key: string, value: string): Promise<ShellResult> {
    const command = `${cmclPath} config ${key} ${value}`;
    return await ShellAPI.execute(command);
  },

  async installMod(url: string): Promise<ShellResult> {
    const command = `${cmclPath} mod --url=${url}`;
    return await ShellAPI.execute(command);
  },

  async searchMod(name: string): Promise<ShellResult> {
    const command = `${cmclPath} mod --info --name=${name}`;
    return await ShellAPI.execute(command);
  },

  async installModpack(url: string): Promise<ShellResult> {
    const command = `${cmclPath} modpack --url=${url}`;
    return await ShellAPI.execute(command);
  },

  async searchModpack(name: string): Promise<ShellResult> {
    const command = `${cmclPath} modpack --info --name=${name}`;
    return await ShellAPI.execute(command);
  },

  async checkForUpdates(): Promise<ShellResult> {
    const command = `${cmclPath} --check-for-updates`;
    return await ShellAPI.execute(command);
  },

  async getAboutInfo(): Promise<ShellResult> {
    const command = `${cmclPath} --about`;
    return await ShellAPI.execute(command);
  }
};
