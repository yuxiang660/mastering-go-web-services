# 问题记录
* go run xxx.exe 直接退出，但没有错误
    - 可能原因：端口被其他程序占用
    - windows端解决办法<br>
        - 检查端口号被占用的PID<br>
        > netstat -aon | findstr <port_number>
        - 查看此PID的程序<br>
        > tasklist | findstr <PID>
        - 杀死此程序<br>
        > taskkill /F /PID <PID>
