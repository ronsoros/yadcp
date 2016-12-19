# yadcp

this adds simple federation to irc:

join any channel on a remote server by doing /join #channel@irc.server.net!port e.g. /join #openstar@home.irc.openstar.pw!6667

## setup

    go get github.com/ronsoros/yadcp
    
## usage

    yadcp [servername] "MOTD String" port
    
    or with ssl:
    yadcp [servername] "MOTD String" port ssl-certificate ssl-key sslport
    
## quality

this is my first Go project and i have not used Go before this program. I know the code is ugly but i made it work ok, now don't be to harsh.
