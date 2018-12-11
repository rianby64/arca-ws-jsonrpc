'use strict';
const conn = new WebSocket("ws://" + document.location.host + "/ws");

conn.onmessage = (e) => {
    const data = JSON.parse(e.data);
    console.log(data);
}

conn.onopen = () => {
    conn.send(JSON.stringify({
        Method: 'read',
        ID: 'my-id',
        Context: {
            source: "Goods"
        }
    }));
}