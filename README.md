# tourl
批量处理 URL/IP/DOMAIN 为 指定形式的URL或IP:PORT形式，懒癌患者福音

# 用法
`cat ips.txt|tourl`

<img width="374" alt="Snipaste_2022-05-06_19-04-52" src="https://user-images.githubusercontent.com/48342077/167122853-1dafd370-fac8-41b9-9b0e-e5b882092bc8.png">


对指定ip文件处理，移除port

`cat ips.txt|tourl    -p 0`

<img width="499" alt="Snipaste_2022-05-06_19-15-05" src="https://user-images.githubusercontent.com/48342077/167122888-5658d774-065c-4a4e-8d13-f6d8d31214c7.png">


对指定ip文件处理，附加常见WEB端口，静默输出

`cat ips.txt|tourl    -P  -q`

<img width="515" alt="Snipaste_2022-05-06_19-13-40" src="https://user-images.githubusercontent.com/48342077/167122909-5cdd6e29-ac16-4dbe-bb7b-0b348c8d83d3.png">


对输入的URL处理，添加端口21，22，附加默认WEB常见端口，扩充为https,指定URI路径，输出到指定文件

`echo  baidu.com:3306/baidu.php|tourl -p 21,22  -ts  -path baidunb.jsp  -P  -o baidu.txt`

<img width="809" alt="Snipaste_2022-05-06_19-09-11" src="https://user-images.githubusercontent.com/48342077/167122959-ee57d047-2922-4cfe-a331-46505317b655.png">


作为 mapcidr和httpx的中转，主动探测C段及title
`mapcidr  -cidr 203.218.154.232/24 |tourl -p 21,22,23  -P   -q |httpx  -sc -title`

<img width="726" alt="Snipaste_2022-05-06_19-22-18" src="https://user-images.githubusercontent.com/48342077/167122987-08da4fae-6bf6-4849-ab69-6584b51abbbf.png">

# 参数
  -P    在常见WEB端口组的基础上指定port，可单独使用
  
  -l string        url文件路径 (default "ips.txt")
        
  -o string        输出到指定文件 (default "result.txt")
        
  -p string        指定port
        
  -path string        URI路径 (default "/")
        
  -q    安静模式，减少输出
  
  -th        转换为HTTP
        
  -ts        转换为HTTPS



