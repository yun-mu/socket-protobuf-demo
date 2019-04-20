# 安全APP


API可视化文档：API


## socket proto 接口

socket 监听地址：<ip:port>
协议：tcp

proto 文件链接：<proto file url>

心跳机制：客户端每隔10s向服务端发送Ping心跳保持在线，即ReqMethod为Ping的请求包
共享地理位置机制：

- 客户端维护一个监听队列，即为要监听的联系人，eg:

  ```java
  // Java

  Queue<String> queue = new LinkedList<String>();
  // socket 处理
  // ...
  // 处理 ResMethod为AddListenPeople 的服务端socket请求
  // 获取userID
  boolean isExist = false;
  for(String q : queue){
      if (q == userID) {
          isExist = true
          break;
      }
  }
  if (!isExist) {
      queue.offer(userID);
  }
  ```

- 客户端维护一个定时器，比如每隔5s拉取一次监听队列中的联系人的地理位置和电量