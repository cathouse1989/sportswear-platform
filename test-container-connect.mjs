// 容器内连通性测试
fetch('http://backend:8080/api/v1/public/products?page=1&pageSize=3&lang=zh')
  .then(r => r.json())
  .then(d => console.log('SUCCESS items:', d && d.data && d.data.items ? d.data.items.length : 0))
  .catch(e => console.log('ERR:', e.message))