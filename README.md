# Falconer: one terminal manager

A hepler tool to manage many remote ssh server.

when you open a native terminal on unix system , Falconer can list all remote ssh server that you have configured as a
pretty table. Then you can choose one to login.

You will get more information about every remote ssh server depend on what's your describe for the server

# Why

I just need a gadget for my work . In my work, i need to manage many remote server machines that deployed some
application, such as java application,mysql and so on.
when i decided to connect some machines, I spend a lot of time checking which applications are deployed on each server
it's making me miserable. Tabby has wonderful **Profiles & connections** feature, but there no space to remark more
information
for every machine or connection. I still need another document to manage information that can not hold with Tabby.

# Install

git clone codes to your local then run

```shell
go install
```

# Configuration 

you can see [example-config.yaml](example-config.yaml). Please move configured file to `$HOME/.ssh/stabby-config.yaml`