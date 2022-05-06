# tourl
批量处理 URL/IP/DOMAIN 为 指定形式的URL或IP:PORT形式，懒癌患者福音

# 用法
cat ips.txt|tourl

![image.png](https://note.youdao.com/yws/res/5/WEBRESOURCE873369acad1e004fdec2b4b8b63658b5)


对指定ip文件处理，移除port
`cat ips.txt|tourl    -p 0`

![image.png](https://note.youdao.com/yws/res/b/WEBRESOURCE6a9dcfab25de4fdc3ec45ada2483a79b)


对指定ip文件处理，附加常见WEB端口，静默输出

`cat ips.txt|tourl    -P  -q`

![image.png](https://note.youdao.com/yws/res/4/WEBRESOURCE1df568e811aed80186574ebb360db964)


对输入的URL处理，添加端口21，22，附加默认WEB常见端口，扩充为https,指定URI路径，输出到指定文件

`echo  baidu.com:3306/baidu.php|tourl -p 21,22  -ts  -path baidunb.jsp  -P  -o baidu.txt`

![image.png](https://note.youdao.com/yws/res/c/WEBRESOURCEd94ac692497b4c0c3503cd6511b4dc6c)


作为 mapcidr和httpx的中转，主动探测C段及title
`mapcidr  -cidr 203.218.154.232/24 |tourl -p 21,22,23  -P   -q |httpx  -sc -title`

![image.png](https://note.youdao.com/yws/res/a/WEBRESOURCE0840a39c3520eace3347c3640d1081ba)

# 参数
  -P    在常见WEB端口组的基础上指定port，可单独使用
  
  -l string        url文件路径 (default "ips.txt")
        
  -o string        输出到指定文件 (default "result.txt")
        
  -p string        指定port
        
  -path string        URI路径 (default "/")
        
  -q    安静模式，减少输出
  
  -th        转换为HTTP
        
  -ts        转换为HTTPS



