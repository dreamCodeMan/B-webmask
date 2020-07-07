# B-webmask
B站防挡弹幕蒙版 .webmask 文件格式解压

文件可分为头部和内容两部分
头部格式

name | desc | type | bytes | offset
---- | --- | --- | --- | ---
tag | 文件类型标识符，目前固定为 MASK, 4个字节 |bytes | 4 | 0
version | 版本号，目前为 1 | Int32 | 4 | 4
check code | 校验码？目前为 2 | Int8 | 1 | 8
segments | 所包含段数，每段时间10秒左右, 一个 3 分钟的视频大概会分成 18 段 | Int32 | 4 | 12
segments meta | 段的元数据，每段包含 16 个字节, 前 8 个字节表示时间，后 8 个字节表示数据 offset | bytes | 16 * segments | 16

B 站 .webmask 文件加载流程：

1. 通过 range 头进行分段加载；首先加载前 0-16 字节(理论上只需读取 0-15 前 16 字节即可)，校验文件类型
2. 校验成功，加载 16-(segments * 16) 字节元数据并进行解析
3. 下载后续最多不超过 22 段的数据，通过 pako 进行解压，根据当前视频播放时间从中选择蒙版进行渲染
4. 大部分视频不会超过 22 段(22 * 10 / 60 约 3 分钟)，超过的话，会按需继续加载后面的数据(每次最多不超过 22 段)


例子：
超过 22 段
// https://www.bilibili.com/video/av21101827

不超过 22 段
// https://www.bilibili.com/video/av46459801
