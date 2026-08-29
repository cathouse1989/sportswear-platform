// /uploads/** 同源反代：门户里轮播图/内容若配置了相对路径图片（如 /uploads/xxx），
// 浏览器端请求本站 /uploads/**，由 Nitro 转发到后端静态文件服务（8080 /uploads），
// 与后端 API 反代同理，彻底避免跨域且无需对外暴露后端端口。
export default defineEventHandler((event) => {
  const config = useRuntimeConfig(event)
  // apiServer 形如 http://backend:8080/api/v1，剥掉路径后缀得到源，再拼上原请求路径（含 /uploads 前缀）
  const origin = (config.apiServer || 'http://localhost:8080/api/v1').replace(/\/api\/v1\/?$/, '')
  return proxyRequest(event, origin + event.path)
})
