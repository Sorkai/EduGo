// 前端配置文件

interface Config {
  // API基础URL
  apiBaseUrl: string;
  // 其他配置项可以在这里添加
}

// 默认配置
const defaultConfig: Config = {
  apiBaseUrl: '/api/v1', // 默认API地址，适用于开发环境
};

// 根据环境变量加载配置
const loadConfig = (): Config => {
  // 无论是开发环境还是生产环境，都优先使用环境变量中的配置
  return {
    ...defaultConfig,
    // 如果有环境变量定义的API地址，则使用环境变量的值
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL || defaultConfig.apiBaseUrl,
  };
};

// 导出配置
const config = loadConfig();
export default config;
