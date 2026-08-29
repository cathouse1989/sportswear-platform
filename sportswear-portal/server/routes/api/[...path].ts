// 后端 API 同源反代（Nitro server 路由，运行时读取 NUXT_API_SERVER，无需构建期注入）。
// 浏览器端统一请求本站 /api/v1/**，由 Nitro 转发到后端服务：
// 局域网/公网（内网穿透）访问都只需暴露 portal 一个地址，
// 后端 8080 端口无需对外开放，同时彻底避免 CORS 问题。
export default defineEventHandler((event) => {
  const config = useRuntimeConfig(event)
  // apiServer 形如 http://backend:8080/api/v1，剥掉路径后缀得到源，
  // 再拼上完整请求路径（含 /api/v1 前缀与查询串）
  const origin = (config.apiServer || 'http://localhost:8080/api/v1').replace(/\/api\/v1\/?$/, '')
  return proxyRequest(event, origin + event.path)
})
