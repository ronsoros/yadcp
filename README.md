# yadcp

this adds simple federation to irc:

join any channel on a remote server by doing /join #channel@irc.server.net!port e.g. /join #openstar@home.irc.openstar.pw!6667.
all joins to #channel are redirected to the remote one by default afterwards.

## setup

    go get github.com/ronsoros/yadcp
    
## usage

    yadcp [servername] "MOTD String" port
    
    or with ssl:
    yadcp [servername] "MOTD String" port ssl-certificate ssl-key sslport
    
## quality

this is my first Go project and i have not used Go before this program. I know the code is ugly but i made it work ok, now don't be to harsh.

## test servers

Connect to irc.openstar.pw or irc.home.openstar.pw and join #openstar

## features

channel-message relaying: 100%
channel-topic relaying: 50%
channel-user relaying: 0%
user-user relaying: 0%
stability: 10%
