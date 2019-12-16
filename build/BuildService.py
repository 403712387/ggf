import os
import stat
import copy
import shutil
import time
import sys

currentTime = time.localtime()
strTime = "%d-%02d-%02d %02d:%02d:%02d" % (currentTime.tm_year, currentTime.tm_mon, currentTime.tm_mday, currentTime.tm_hour, currentTime.tm_min,currentTime.tm_sec)

# 服务模块
serviceName = "ggf"

# git信息
gitBranch = "unknown"
gitCommitId = "unknown"

#编译参数，支持debug,race
compileArg = ""

#------------------------函数的定义-------------------------#


#清理
def cleanFiles(path):
    if os.path.exists(path):
        shutil.rmtree(path)

#解析参数
def parseArgs():
    global compileArg

    if "race" in sys.argv:
        compileArg = "-race"

    if "debug" in sys.argv:
        compileArg = '''-gcflags "-N -l"'''

#下载依赖的包
def downloadThirdLibrary():
    librarys = ["github.com/btfak/sntp", "github.com/sirupsen/logrus", "github.com/shirou/gopsutil", "github.com/segmentio/kafka-go", "github.com/mattn/go-sqlite3"]
    for library in librarys:
        os.system("go get %s"%library)

#获取git的信息（获取当前分支以及commit id）
def getGitInfo():
    global gitBranch, gitCommitId

    gitDir = "../.git"

    #获取分支信息
    branchFile = os.path.join(gitDir, "HEAD")
    if os.path.exists(branchFile):
        with open(branchFile, "r") as f:
            line = f.readline()
            line = line.strip()
            splits = line.split("/")
            if len(splits) > 0:
                gitBranch = splits[-1]

    # 获取commit id
    commitIdFile = os.path.join(gitDir + "/refs/heads" , gitBranch)
    if os.path.exists(commitIdFile):
        with open(commitIdFile) as f:
            line = f.readline()
            line = line.strip()
            gitCommitId = line

#编译各个模块
def compileService():
    global serviceName, compileArg, gitBranch, gitCommitId

    compileSuccessful = False

    # 切换目录
    currentPath = os.getcwd()
    os.chdir("../src")

    # 格式话git信息
    git = "-X HostServiceModule.GitBranch=%s -X HostServiceModule.GitCommitID=%s"%(gitBranch, gitCommitId)


    # 获取当前GOPATH
    currentGoPath = ""
    pipe = os.popen("go env GOPATH")
    lines = pipe.readlines()
    if len(lines) > 0 :
        currentGoPath = lines[0].strip("\n")

    # 编译
    projectPath = os.getcwd()[:-4]
    goPathEnv = "export GOPATH=%s:%s"%(currentGoPath,projectPath)
    os.system(goPathEnv + " && go clean")
    os.system(goPathEnv + " && go clean -r")
    os.system(goPathEnv + " && go clean -cache")
    compile = '''go build -ldflags "%s" %s -o ./bin/%s/%s main.go'''%(git, compileArg, serviceName, serviceName)
    print(goPathEnv + " && " + compile)
    if os.system(goPathEnv + " && " + compile) == 0:
        compileSuccessful = True

    os.chdir(currentPath)
    return compileSuccessful

# 拷贝配置文件
def copyConfigFile():
    global serviceName

    src = "../config"
    dst = "../src/bin/%s"%serviceName
    copyFiles(src, dst)

    # 配置文件要放在config目录下
    src = "../src/bin/%s/config.json"%serviceName
    dst = "../src/bin/%s/config"%serviceName
    copyFiles(src, dst)
    os.remove(src)

#修改文件的权限
def processFilePromission(path):
    files = os.listdir(path)
    for file in files:
        fileName = os.path.join(path, file)
        if not os.path.isfile(fileName):
            continue

        #对于sh结束的文件，修改权限
        if fileName.endswith(".sh"):
            os.chmod(fileName, stat.S_IRWXU | stat.S_IRWXG | stat.S_IRWXO)

#拷贝文件或者目录
def copyFiles(source, destination):

    #复制文件(要注意权限和软连接这种情况)
    def copyFile(sourceFile, destDir):
        if not os.path.exists(sourceFile):
            return

        if not os.path.exists(destDir):
            os.makedirs(destDir)

        if os.path.islink(sourceFile):  #复制软连接
            currentPath = os.getcwd()
            symbolLink = os.readlink(sourceFile)
            os.chdir(destDir)
            os.symlink(symbolLink, os.path.basename(sourceFile))
            os.chdir(currentPath)
        elif os.path.isfile(sourceFile):    #复制文件
            with open(sourceFile, "rb") as input:
                with open(os.path.join(destDir, os.path.basename(sourceFile)), "wb") as output:
                    output.write(input.read())
        os.chmod(os.path.join(destDir, os.path.basename(sourceFile)), os.stat(sourceFile).st_mode)

    if not os.path.exists(source):
        print("copy %s to %s fail, not find %s"%(source, destination, source))
        return

    # 目标文件夹一定要存在
    if not os.path.exists(destination):
        os.makedirs(destination)

    if os.path.isdir(source):   #复制整个目录下的文件
        for path, directorys, files in os.walk(source):
            subPath = path[len(source): ]

            # 创建目录
            if subPath.startswith("/"):
                subPath = subPath[1:]
            destinationPath = os.path.join(destination, subPath)
            if not os.path.exists(destinationPath):
                os.makedirs(destinationPath)

            # 复制目录下中的文件
            for file in files:
                copyFile(os.path.join(path, file), destinationPath)
    elif os.path.isfile(source):    # 复制单个文件
        copyFile(source, destination)


#修改脚本中的结束符，\r\n换为\n
def formatLineBrak():
    global serviceName

    fileNames = ["../src/bin/" + serviceName + "/start.sh", "../src/bin/" + serviceName + "/stop.sh"]
    for fileName in fileNames:
        if not os.path.exists(fileName):
            continue

        fileData = ""
        with open(fileName, "r")  as file:
            for lineData in file:
                lineData = lineData.replace("\r\n", "\n")
                fileData += lineData

            # 向启动脚本中写入内容
        with open(fileName, "w")  as file:
            file.write(fileData)

#构建服务
def buildService():
    global serviceName
    outputDir = "../src/bin/" + serviceName
    serviceDir = "./" + serviceName
    parseArgs()
    downloadThirdLibrary()
    cleanFiles(outputDir)
    cleanFiles(serviceDir)
    getGitInfo()

    #编译各个模块
    if not compileService():
        print("\n--------------compile fail at %s--------------" % (strTime))
        return -1

    #拷贝文件
    copyConfigFile()

    #处理脚本
    formatLineBrak()

    #修改文件的权限
    processFilePromission(outputDir)

    #移动到当前目录
    print("move dir %s to %s"%(outputDir, serviceDir))
    copyFiles(outputDir, serviceDir)
    print("\n--------------compile successful at %s--------------"%(strTime))
    return 0

#------------------------函数的调用-------------------------#
buildService()