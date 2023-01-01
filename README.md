## 接口设计：
ffmpeg裁剪影片，将影片裁剪成指定时间段的影片，返回给客户端。

## 业务逻辑分解：
设计一个处理带有影片URL地址、起始时间、终止时间信息请求的接口，接口返回裁剪后的影片URL地址。

1. 配置好应用所需要的参数，使得包括数据库、日志器、输出文件夹等可以正常使用
2. 记录影片URL地址、起始时间、终止时间、请求时间信息，将数据持久化到数据库
3. 从URL中获取影片文件，下载文件函数设计
4. 用封装好的生成器ffmpeg-go对影片根据需求进行处理
5. 当ffmpeg处理完成后，改写数据库中记录的id_done状态

## 使用方法：
1. 修改configs/conf.yml配置文件
2. 启动服务
3. 提交请求：用json格式的数据POST方式请求接口 http://localhost:port/api/v1/clip
4. 查看是否完成：GET方式请求接口 http://localhost:port/api/v1/clip?id=1 id为数据库中记录的id

参考请求json:

```json
{
    "video_url": "https://cdn.coverr.co/videos/coverr-a-vintage-mercedes-benz-w123-at-a-motor-show-7241/1080p.mp4",
    "start_time": "1",
    "end_time": "3"
}
```

## 未完成部分的思路:
1. 接口鉴别用户身份：使用JWT, 通过用户的token来鉴别用户身份。我在另一个项目实现过
2. 接口支持主动推送剪辑完成事件给用户