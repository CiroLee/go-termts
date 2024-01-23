# go-termts           
> common and useful terminal commands written in go

## Install
```shell
go install github.com/CiroLee/go-termts
```

## Usage 
```shell
go-termts <command> [flags]
```

## Features     

### license     
output LICENSE in an interactive way         

### commit        
shortcut for git commit -m, support zh(for Chinese) and en(for English) flags       

### repo      
open current git project in your default browser     

### config        
download common used config files, support `prettier`, `commitlint`, `vscode`(vscode-settings)

### ip        
output the local ip

### alias        
output alias from your `.zshrc` file

### dgit      
download github repository. support custom branch and custom path. default branch is repo's default branch and default path is current directory.     

```shell
# help
go-termts dgit <repo> [--branch] [--dst]
# example
go-termts dgit https://github.com/CiroLee/go-termts --branch=main --dst=demo
```