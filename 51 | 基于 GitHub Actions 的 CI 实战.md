在 Go 项目开发中，我们要频繁地执行静态代码检查、测试、编译、构建等操作。如果每一步我们都手动执行，效率低不说，还容易出错。所以，我们通常借助 CI 系统来自动化执行这些操作。

当前业界有很多优秀的 CI 系统可供选择，例如 CircleCI、TravisCI、Jenkins、CODING、GitHub Actions 等。这些系统在设计上大同小异，为了减少你的学习成本，我选择了相对来说容易实践的 GitHub Actions，来给你展示如何通过 CI 来让工作自动化。

**GitHub Actions 的基本用法**

GitHub Actions 的基本用法GitHub Actions 是 GitHub 为托管在 github.com 站点的项目提供的持续集成服务，于 2018 年 10 月推出。

GitHub Actions 具有以下功能特性：

- 提供原子的 actions 配置和组合 actions 的 workflow 配置两种能力。
- 全局配置基于YAML 配置，兼容主流 CI/CD 工具配置。
- Actions/Workflows 基于事件触发，包括 Event restrictions、Webhook events、Scheduled events、External events。
- 提供可供运行的托管容器服务，包括 Docker、VM，可运行 Linux、macOS、Windows 主流系统。
- 提供主流语言的支持，包括 Node.js、Python、Java、Ruby、PHP、Go、Rust、.NET。
- 提供实时日志流程，方便调试。
- 提供平台内置的 Actions与第三方提供的 Actions，开箱即用。

**GitHub Actions 的基本概念**

在构建持续集成任务时，我们会在任务中心完成各种操作，比如克隆代码、编译代码、运行单元测试、构建和发布镜像等。GitHub 把这些操作称为 Actions。

Actions 在很多项目中是可以共享的，GitHub 允许开发者将这些可共享的 Actions 上传到GitHub 的官方 Actions 市场，开发者在 Actions 市场中可以搜索到他人提交的 Actions。另外，还有一个 awesome actions 的仓库，里面也有不少的 Action 可供开发者使用。如果你需要某个 Action，不必自己写复杂的脚本，直接引用他人写好的 Action 即可。整个持续集成过程，就变成了一个 Actions 的组合。

Action 其实是一个独立的脚本，可以将 Action 存放在 GitHub 代码仓库中，通过<userName>/<repoName>的语法引用 Action。例如，actions/checkout@v2表示https://github.com/actions/checkout这个仓库，tag 是 v2。actions/checkout@v2也代表一个 Action，作用是安装 Go 编译环境。GitHub 官方的 Actions 都放在 github.com/actions 里面。

GitHub Actions 有一些自己的术语，下面我来介绍下。

- workflow（工作流程）：一个  .yml  文件对应一个 workflow，也就是一次持续集成。一个 GitHub 仓库可以包含多个 workflow，只要是在  .github/workflow  目录下的  .yml  文件都会被 GitHub 执行。
- job（任务）：一个 workflow 由一个或多个 job 构成，每个 job 代表一个持续集成任务。
- step（步骤）：每个 job 由多个 step 构成，一步步完成。
- action（动作）：每个 step 可以依次执行一个或多个命令（action）。
- on：一个 workflow 的触发条件，决定了当前的 workflow 在什么时候被执行。

**workflow 文件介绍**

GitHub Actions 配置文件存放在代码仓库的.github/workflows目录下，文件后缀为.yml，支持创建多个文件，文件名可以任意取，比如iam.yml。GitHub 只要发现.github/workflows目录里面有.yml文件，就会自动运行该文件，如果运行过程中存在问题，会以邮件的形式通知到你。

workflow 文件的配置字段非常多，如果你想详细了解，可以查看官方文档。这里，我来介绍一些基本的配置字段。

name

name字段是 workflow 的名称。如果省略该字段，默认为当前 workflow 的文件名。

```yml
name: GitHub Actions Demo
```

on

on字段指定触发 workflow 的条件，通常是某些事件。

```yml
on: push
```

上面的配置意思是，push事件触发 workflow。on字段也可以是事件的数组，例如:

```yml
on: [push, pull_request]
```

上面的配置意思是，push事件或pull_request事件都可以触发 workflow。

想了解完整的事件列表，你可以查看官方文档[https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows ]。除了代码库事件，GitHub Actions 也支持外部事件触发，或者定时运行。

on.<push|pull_request>.<tags|branches>

指定触发事件时，我们可以限定分支或标签。

```yml
on:
  push:
    branches:
      - master
```

上面的配置指定，只有master分支发生push事件时，才会触发 workflow。

jobs.<job_id>.name

workflow 文件的主体是jobs字段，表示要执行的一项或多项任务。

jobs字段里面，需要写出每一项任务的job_id，具体名称自定义。job_id里面的name字段是任务的说明。

```yml
jobs:
  my_first_job:
    name: My first job
  my_second_job:
    name: My second job
```

上面的代码中，jobs字段包含两项任务，job_id分别是my_first_job和my_second_job。

jobs.<job_id>.needs

needs字段指定当前任务的依赖关系，即运行顺序。

```yml
jobs:
  job1:
  job2:
    needs: job1
  job3:
    needs: [job1, job2]
```

上面的代码中，job1必须先于job2完成，而job3等待job1和job2完成后才能运行。因此，这个 workflow 的运行顺序为：job1、job2、job3。

jobs.<job_id>.runs-on

runs-on字段指定运行所需要的虚拟机环境，它是必填字段。目前可用的虚拟机如下：

- ubuntu-latest、ubuntu-18.04 或 ubuntu-16.04。
- windows-latest、windows-2019 或 windows-2016。
- macOS-latest 或 macOS-10.14。

下面的配置指定虚拟机环境为ubuntu-18.04。

```yml
runs-on: ubuntu-18.04
```

jobs.<job_id>.steps

steps字段指定每个 Job 的运行步骤，可以包含一个或多个步骤。每个步骤都可以指定下面三个字段。

- jobs.<job_id>.steps.name：步骤名称。
- jobs.<job_id>.steps.run：该步骤运行的命令或者 action。
- jobs.<job_id>.steps.env：该步骤所需的环境变量。

下面是一个完整的 workflow 文件的范例：

```yml
name: Greeting from Mona
on: push

jobs:
  my-job:
    name: My Job
    runs-on: ubuntu-latest
    steps:
    - name: Print a greeting
      env:
        MY_VAR: Hello! My name is
        FIRST_NAME: Lingfei
        LAST_NAME: Kong
      run: |
        echo $MY_VAR $FIRST_NAME $LAST_NAME.
```

上面的代码中，steps字段只包括一个步骤。该步骤先注入三个环境变量，然后执行一条 Bash 命令。

uses

uses 可以引用别人已经创建的 actions，就是上面说的 actions 市场中的 actions。引用格式为userName/repoName@verison，例如uses: actions/setup-go@v1。

with

with 指定 actions 的输入参数。每个输入参数都是一个键 / 值对。输入参数被设置为环境变量，该变量的前缀为 INPUT_，并转换为大写。

这里举个例子：我们定义 hello_world 操作所定义的三个输入参数（first_name、middle_name 和 last_name），这些输入变量将被 hello-world 操作作为 INPUT_FIRST_NAME、INPUT_MIDDLE_NAME 和 INPUT_LAST_NAME 环境变量使用。

```yml
jobs:
  my_first_job:
    steps:
      - name: My first step
        uses: actions/hello_world@master
        with:
          first_name: Lingfei
          middle_name: Go
          last_name: Kong
```

run

run指定执行的命令。可以有多个命令，例如：

```yml
- name: Build
      run: |
      go mod tidy
      go build -v -o helloci .
```

id

id是 step 的唯一标识。

**GitHub Actions 的进阶用法**

上面，我介绍了 GitHub Actions 的一些基本知识，这里我再介绍下 GitHub Actions 的进阶用法。

**为工作流加一个 Badge**

在 action 的面板中，点击Create status badge就可以复制 Badge 的 Markdown 内容到 README.md 中。

之后，我们就可以直接在 README.md 中看到当前的构建结果：

![img](https://static001.geekbang.org/resource/image/45/af/453a97b0776281873dee5671c53347af.png?wh=1280x765)

**使用构建矩阵**

如果我们想在多个系统或者多个语言版本上测试构建，就需要设置构建矩阵。例如，我们想在多个操作系统、多个 Go 版本下跑测试，可以使用如下 workflow 配置：

```yml
name: Go Test

on: [push, pull_request]

jobs:

  helloci-build:
    name: Test with go ${{ matrix.go_version }} on ${{ matrix.os }}
    runs-on: ${{ matrix.os }}

    strategy:
      matrix:
        go_version: [1.15, 1.16]
        os: [ubuntu-latest, macOS-latest]

    steps:

      - name: Set up Go ${{ matrix.go_version }}
        uses: actions/setup-go@v2
        with:
          go-version: ${{ matrix.go_version }}
        id: go
```

**使用 Secrets**

在构建过程中，我们可能需要用到ssh或者token等敏感数据，而我们不希望这些数据直接暴露在仓库中，此时就可以使用secrets。

我们在对应项目中选择Settings-> Secrets，就可以创建secret，如下图所示：

![img](https://static001.geekbang.org/resource/image/c0/d3/c00b11a1709838c1a205ace7976768d3.png?wh=1920x1046)

配置文件中的使用方法如下：

```yml
name: Go Test
on: [push, pull_request]
jobs:
  helloci-build:
    name: Test with go
    runs-on: [ubuntu-latest]
    environment:
      name: helloci
    steps:
      - name: use secrets
        env:
          super_secret: ${{ secrets.YourSecrets }}
```

secret name 不区分大小写，所以如果新建 secret 的名字是 name，使用时用 secrets.name 或者 secrets.Name 都是可以的。而且，就算此时直接使用 echo 打印 secret , 控制台也只会打印出*来保护 secret。

这里要注意，你的 secret 是属于某一个环境变量的，所以要指明环境的名字：environment.name。上面的 workflow 配置中的secrets.YourSecrets属于helloci环境。

**使用 Artifact 保存构建产物**

在构建过程中，我们可能需要输出一些构建产物，比如日志文件、测试结果等。这些产物可以使用 Github Actions Artifact 来存储。你可以使用action/upload-artifact 和 download-artifact 进行构建参数的相关操作。

这里我以输出 Jest 测试报告为例来演示下如何保存 Artifact 产物。Jest 测试后的测试产物是 coverage：

```yml
steps:
      - run: npm ci
      - run: npm test

      - name: Collect Test Coverage File
        uses: actions/upload-artifact@v1.0.0
        with:
          name: coverage-output
          path: coverage
```

行成功后，我们就能在对应 action 面板看到生成的 Artifact：

![img](https://static001.geekbang.org/resource/image/4c/66/4c4a8d6aec12a5dd1cdc80d238472566.png?wh=1280x208)

**GitHub Actions 实战**

上面，我介绍了 GitHub Actions 的用法，接下来我们就来实战下，看下使用 GitHub Actions 的 6 个具体步骤。

第一步，创建一个测试仓库。

登陆GitHub 官网，点击 New repository 创建，如下图所示：

![img](https://static001.geekbang.org/resource/image/6d/a0/6d76d02f0418671a32f5346fccf616a0.png?wh=1920x810)

这里，我们创建了一个叫helloci的测试项目。

第二步，将新的仓库 clone 下来，并添加一些文件：

```shell
$ git clone https://github.com/marmotedu/helloci
```

你可以克隆marmotedu/helloci，并将里面的文件拷贝到你创建的项目仓库中。

第三步，创建 GitHub Actions workflow 配置目录：

```shell
$ mkdir -p .github/workflows                     
```

第四步，创建 GitHub Actions workflow 配置。

在.github/workflows目录下新建helloci.yml文件，内容如下：

```yml
name: Go Test

on: [push, pull_request]

jobs:

  helloci-build:
    name: Test with go ${{ matrix.go_version }} on ${{ matrix.os }}
    runs-on: ${{ matrix.os }}
    environment:
      name: helloci

    strategy:
      matrix:
        go_version: [1.16]
        os: [ubuntu-latest]

    steps:

      - name: Set up Go ${{ matrix.go_version }}
        uses: actions/setup-go@v2
        with:
          go-version: ${{ matrix.go_version }}
        id: go

      - name: Check out code into the Go module directory
        uses: actions/checkout@v2

      - name: Tidy
        run: |
          go mod tidy

      - name: Build
        run: |
          go build -v -o helloci .

      - name: Collect main.go file
        uses: actions/upload-artifact@v1.0.0
        with:
          name: main-output
          path: main.go

      - name: Publish to Registry
        uses: elgohr/Publish-Docker-GitHub-Action@master
        with:
          name: ccr.ccs.tencentyun.com/marmotedu/helloci:beta  # docker image 的名字
          username: ${{ secrets.DOCKER_USERNAME}} # 用户名
          password: ${{ secrets.DOCKER_PASSWORD }} # 密码
          registry: ccr.ccs.tencentyun.com # 腾讯云Registry
          dockerfile: Dockerfile # 指定 Dockerfile 的位置
          tag_names: true # 是否将 release 的 tag 作为 docker image 的 tag
```

上面的 workflow 文件定义了当 GitHub 仓库有push、pull_request事件发生时，会触发 GitHub Actions 工作流程，流程中定义了一个任务（Job）helloci-build，Job 中包含了多个步骤（Step），每个步骤又包含一些动作（Action）。

上面的 workflow 配置会按顺序执行下面的 6 个步骤。

1. 准备一个 Go 编译环境。
2. 从marmotedu/helloci下载源码。
3. 添加或删除缺失的依赖包。
4. 编译 Go 源码。
5. 上传构建产物。
6. 构建镜像，并将镜像 push 到ccr.ccs.tencentyun.com/marmotedu/helloci:beta

第五步，在 push 代码之前，我们需要先创建DOCKER_USERNAME和DOCKER_PASSWORD secret。

其中，DOCKER_USERNAME保存腾讯云镜像服务（CCR）的用户名，DOCKER_PASSWORD保存 CCR 的密码。我们将这两个 secret 保存在helloci Environments 中，如下图所示：

![img](https://static001.geekbang.org/resource/image/c0/d3/c00b11a1709838c1a205ace7976768d3.png?wh=1920x1046)

第六步，将项目 push 到 GitHub，触发 workflow 工作流：

```shell
$ git add .
$ git push origin master
```

打开我们的仓库 Actions 标签页，可以发现 GitHub Actions workflow 正在执行：

![img](https://static001.geekbang.org/resource/image/1a/8a/1afb7860d68635c5e3eaba4ff8da208a.png?wh=1920x691)

然后，选择其中一个构建记录，查看其运行详情（具体可参考chore: update step name Go Test #10）：

![img](https://static001.geekbang.org/resource/image/48/4f/481f64aabccf30ed61d0a7c85ab30d4f.png?wh=1920x1084)

你可以看到，Go Test工作流程执行了 6 个 Job，每个 Job 执行了下面这些自定义 Step：

1. Set up Go 1.16。
2. Check out code into the Go module directory。
3. Tidy。
4. Build。
5. Collect main.go file。
6. Publish to Registry。

其他步骤是 GitHub Actions 自己添加的步骤：Setup Job、Post Check out code into the Go module directory、Complete job。点击每一个步骤，你都能看到它们的详细输出。

**IAM GitHub Actions 实战**

接下来，我们再来看下 IAM 项目的 GitHub Actions 实战。

假设 IAM 项目根目录为 ${IAM_ROOT}，它的 workflow 配置文件为：

```yml
$ cat ${IAM_ROOT}/.github/workflows/iamci.yaml
name: IamCI

on:
  push:
    branchs:
    - '*'
  pull_request:
    types: [opened, reopened]

jobs:

  iamci:
    name: Test with go ${{ matrix.go_version }} on ${{ matrix.os }}
    runs-on: ${{ matrix.os }}
    environment:
      name: iamci

    strategy:
      matrix:
        go_version: [1.16]
        os: [ubuntu-latest]

    steps:

      - name: Set up Go ${{ matrix.go_version }}
        uses: actions/setup-go@v2
        with:
          go-version: ${{ matrix.go_version }}
        id: go

      - name: Check out code into the Go module directory
        uses: actions/checkout@v2

      - name: Run go modules Tidy
        run: |
          make tidy

      - name: Generate all necessary files, such as error code files
        run: |
          make gen

      - name: Check syntax and styling of go sources
        run: |
          make lint

      - name: Run unit test and get test coverage
        run: |
          make cover

      - name: Build source code for host platform
        run: |
          make build

      - name: Collect Test Coverage File
        uses: actions/upload-artifact@v1.0.0
        with:
          name: main-output
          path: _output/coverage.out

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v1

      - name: Login to DockerHub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build docker images for host arch and push images to registry
        run: |
          make push
```

上面的 workflow 依次执行了以下步骤：

1. 设置 Go 编译环境。
2. 下载 IAM 项目源码。
3. 添加 / 删除不需要的 Go 包。
4. 生成所有的代码文件。
5. 对 IAM 源码进行静态代码检查。
6. 运行单元测试用例，并计算单元测试覆盖率是否达标。
7. 编译代码。
8. 收集构建产物_output/coverage.out。
9. 配置 Docker 构建环境。
10. 登陆 DockerHub。
11. 构建 Docker 镜像，并 push 到 DockerHub。

IamCI workflow 运行历史如下图所示：

![img](https://static001.geekbang.org/resource/image/2b/b0/2b542f9101be0c3a83576fb99bf882b0.png?wh=1920x844)

IamCI workflow 的其中一次工作流程运行结果如下图所示：

![img](https://static001.geekbang.org/resource/image/e9/6a/e9ebf13fdb6e4f41a1b00406e646ec6a.png?wh=1920x887)

**总结**

在 Go 项目开发中，我们需要通过 CI 任务来将需要频繁操作的任务自动化，这不仅可以提高开发效率，还能减少手动操作带来的失误。这一讲，我选择了最易实践的 GitHub Actions，来给你演示如何构建 CI 任务。

GitHub Actions 支持通过 push 事件来触发 CI 流程。一个 CI 流程其实就是一个 workflow，workflow 中包含多个任务，这些任务是可以并行执行的。一个任务又包含多个步骤，每一步又由多个动作组成。动作（Action）其实是一个命令 / 脚本，用来完成我们指定的任务，如编译等。

因为 GitHub Actions 内容比较多，这一讲只介绍了一些核心的知识，更详细的 GitHub Actions 教程，你可以参考 官方中文文档。

**课后练习**

使用 CODING 实现 IAM 的 CI 任务，并思考下：GitHub Actions 和 CODING 在 CI 任务构建上，有没有本质的差异？

这一讲，我们借助 GitHub Actions 实现了 CI，请你结合前面所学的知识，实现 IAM 的 CD 功能。欢迎提交 Pull Request