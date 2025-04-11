// 关闭连接
const closeConnection = (proxy: any) => function (id: any) {
    proxy.$http.delete('/connections/' + id);
}

// 关闭所有连接
const closeAllConnection = (proxy: any) => function () {
    proxy.$http.delete('/connections');
}


export default function createConnApi(proxy: any) {
    return {
        closeConnection: closeConnection(proxy),
        closeAllConnection: closeAllConnection(proxy),
    }
}