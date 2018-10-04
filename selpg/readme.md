## 设计说明
该实用程序从标准输入或从作为命令行参数给出的文件名读取文本输入。它允许用户指定来自该输入并随后将被输出的页面范围。例如，如果输入含有 100 页，则用户可指定只打印第 35 至 65 页。这种特性有实际价值，因为在打印机上打印选定的页面避免了浪费纸张。另一个示例是，原始文件很大而且以前已打印过，但某些页面由于打印机卡住或其它原因而没有被正确打印。在这样的情况下，则可用该工具来只打印需要打印的页面。

该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。

selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。

### 参数处理
**“-sNumber”和“-eNumber”强制选项：**
selpg 要求用户用两个命令行参数“-sNumber”（例如，“-s10”表示从第 10 页开始）和“-eNumber”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。selpg 对所给的页号进行合理性检查；换句话说，它会检查两个数字是否为有效的正整数以及结束页是否不小于起始页。这两个选项，“-sNumber”和“-eNumber”是强制性的，而且必须是命令行上在命令名 selpg 之后的头两个参数：
```bash
$ selpg -s10 -e20 ...
```
（... 是命令的余下部分，下面对它们做了描述）。

**“-lNumber”和“-f”可选选项：**
selpg 可以处理两种输入文本：

*类型 1：* 该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-lNumber”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。

选择 72 作为缺省值是因为在行打印机上这是很常见的页长度。这样做的意图是将最常见的命令用法作为缺省值，这样用户就不必输入多余的选项。该缺省值可以用“-lNumber”选项覆盖，如下所示：
```bash
$ selpg -s10 -e20 -l66 ...
```
这表明页有固定长度，每页为 66 行。

*类型 2：* 该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界。该格式与“每页行数固定”格式相比的好处在于，当每页的行数有很大不同而且文件有很多页时，该格式可以节省磁盘空间。在含有文本的行后面，类型 2 的页只需要一个字符 ― 换页 ― 就可以表示该页的结束。打印机会识别换页符并自动根据在新的页开始新行所需的行数移动打印头。

将这一点与类型 1 比较：在类型 1 中，文件必须包含 PAGELEN - CURRENTPAGELEN 个新的行以将文本移至下一页，在这里 PAGELEN 是固定的页大小而 CURRENTPAGELEN 是当前页上实际文本行的数目。在此情况下，为了使打印头移至下一页的页首，打印机实际上必须打印许多新行。这在磁盘空间利用和打印机速度方面效率都很低（尽管实际的区别可能不太大）。

类型 2 格式由“-f”选项表示，如下所示：
```bash
$ selpg -s10 -e20 -f ...
```
该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。

注：“-lNumber”和“-f”选项是互斥的。

**“-dDestination”可选选项：**
selpg 还允许用户使用“-dDestination”选项将选定的页直接发送至打印机。这里，“Destination”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。该目的地应该存在 ― selpg 不检查这一点。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“Destination”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出。这一特性是用 popen() 系统调用实现的，该系统调用允许一个进程打开到另一个进程的管道，将管道用于输出或输入。在下面的示例中，我们打开到命令
```bash
$ lp -dDestination
```
的管道以便输出，并写至该管道而不是标准输出：
```bash
$ selpg -s10 -e20 -dlp1
```
该命令将选定的页作为打印作业发送至 lp1 打印目的地。您应该可以看到类似“request id is lp1-6”的消息。该消息来自 lp 命令；它显示打印作业标识。如果在运行 selpg 命令之后立即运行命令 lpstat -t | grep lp1 ，您应该看见 lp1 队列中的作业。如果在运行 lpstat 命令前耽搁了一些时间，那么您可能看不到该作业，因为它一旦被打印就从队列中消失了。

### 输入处理
一旦处理了所有的命令行参数，就使用这些指定的选项以及输入、输出源和目标来开始输入的实际处理。
selpg 通过以下方法记住当前页号：如果输入是每页行数固定的，则 selpg 统计新行数，直到达到页长度后增加页计数器。如果输入是换页定界的，则 selpg 改为统计换页符。这两种情况下，只要页计数器的值在起始页和结束页之间这一条件保持为真，selpg 就会输出文本（逐行或逐字）。当那个条件为假（也就是说，页计数器的值小于起始页或大于结束页）时，则 selpg 不再写任何输出。

## 使用
为了演示最终用户可以如何应用我们所介绍的一些原则，下面给出了可使用的 selpg 命令字符串示例：
```bash
$ selpg -s1 -e1 input_file
```
该命令将把“input_file”的第 1 页写至标准输出（也就是屏幕），因为这里没有重定向或管道。
```bash
$ selpg -s1 -e1 < input_file
```
该命令与示例 1 所做的工作相同，但在本例中，selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“input_file”而不是显式命名的文件名参数。输入的第 1 页被写至屏幕。
```bash
$ other_command | selpg -s10 -e20
```
“other_command”的标准输出被 shell／内核重定向至 selpg 的标准输入。将第 10 页到第 20 页写至 selpg 的标准输出（屏幕）。
```bash
$ selpg -s10 -e20 input_file >output_file
```
selpg 将第 10 页到第 20 页写至标准输出；标准输出被 shell／内核重定向至“output_file”。
```bash
$ selpg -s10 -e20 input_file 2>error_file
```
selpg 将第 10 页到第 20 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error_file”。请注意：在“2”和“>”之间不能有空格；这是 shell 语法的一部分（请参阅“man bash”或“man sh”）。
```bash
$ selpg -s10 -e20 input_file >output_file 2>error_file
```
selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至“error_file”。当“input_file”很大时可使用这种调用；您不会想坐在那里等着 selpg 完成工作，并且您希望对输出和错误都进行保存。
```bash
$ selpg -s10 -e20 input_file >output_file 2>/dev/null
```
selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。
```bash
$ selpg -s10 -e20 input_file >/dev/null
```
selpg 将第 10 页到第 20 页写至标准输出，标准输出被丢弃；错误消息在屏幕出现。这可作为测试 selpg 的用途，此时您也许只想（对一些测试情况）检查错误消息，而不想看到正常输出。
```bash
$ selpg -s10 -e20 input_file | other_command
```
selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，第 10 页到第 20 页被写至该标准输入。“other_command”的示例可以是 lp，它使输出在系统缺省打印机上打印。“other_command”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“other_command”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。
```bash
$ selpg -s10 -e20 input_file 2>error_file | other_command
```
与上面的示例 9 相似，只有一点不同：错误消息被写至“error_file”。

在以上涉及标准输出或标准错误重定向的任一示例中，用“>>”替代“>”将把输出或错误数据附加在目标文件后面，而不是覆盖目标文件（当目标文件存在时）或创建目标文件（当目标文件不存在时）。

以下所有的示例也都可以（有一个例外）结合上面显示的重定向或管道命令。我没有将这些特性添加到下面的示例，因为我认为它们在上面示例中的出现次数已经足够多了。例外情况是您不能在任何包含“-dDestination”选项的 selpg 调用中使用输出重定向或管道命令。实际上，您仍然可以对标准错误使用重定向或管道命令，但不能对标准输出使用，因为没有任何标准输出 — 正在内部使用 popen() 函数由管道将它输送至 lp 进程。
```bash
$ selpg -s10 -e20 -l66 input_file
```
该命令将页长设置为 66 行，这样 selpg 就可以把输入当作被定界为该长度的页那样处理。第 10 页到第 20 页被写至 selpg 的标准输出（屏幕）。
```bash
$ selpg -s10 -e20 -f input_file
```
假定页由换页符定界。第 10 页到第 20 页被写至 selpg 的标准输出（屏幕）。
```bash
$ selpg -s10 -e20 -dlp1 input_file
```
第 10 页到第 20 页由管道输送至命令“lp -dlp1”，该命令将使输出在打印机 lp1 上打印。

最后一个示例将演示 Linux shell 的另一特性：
```bash
$ selpg -s10 -e20 input_file > output_file 2>error_file &
```
该命令利用了 Linux 的一个强大特性，即：在“后台”运行进程的能力。在这个例子中发生的情况是：“进程标识”（pid）如 1234 将被显示，然后 shell 提示符几乎立刻会出现，使得您能向 shell 输入更多命令。同时，selpg 进程在后台运行，并且标准输出和标准错误都被重定向至文件。这样做的好处是您可以在 selpg 运行时继续做其它工作。

您可以通过运行命令 ps（代表“进程状态”）检查它是否仍在运行或已经完成。该命令会显示数行信息，每行代表一个从该 shell 会话启动的进程（包括 shell 本身）。如果 selpg 仍在运行，您也将看到表示它的一项信息。您也可以用命令“kill -l5 1234”杀死正在运行的 selpg 进程。如果这不起作用，可以尝试用“kill -9 1234”。警告：在对任何重要进程使用该命令前，请阅读“man kill”。

## 测试结果
用 Python 生成测试文件。
```python
f = open("test", "w+")
for i in range(500):
        f.write("line %d\n" % i)
f.close()
```
测试文件`test`：
```bash
line 0
line 1
line 2
line 3
line 4
line 5
line 6
...
ine 497
line 498
line 499
```
测试一：
```bash
$ ./selpg -s 1 -e 1  -l 15 test
line 0
line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
done
```
测试二：
```bash
$ ./selpg -s 2 -e 3 -l 6 test
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
line 15
line 16
line 17
done
```
测试三：
```bash
$ ./selpg -s 1 -e 1 -l 10 < test
line 0
line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
done
```
测试四：
```bash
$ cat test | ./selpg -s 2 -e 3 -l 10
line 10
line 11
line 12
line 13
line 14
line 15
line 16
line 17
line 18
line 19
line 20
line 21
line 22
line 23
line 24
line 25
line 26
line 27
line 28
line 29
done
```
测试五：
```bash
$ ./selpg -s 2 -e 3 -l 10 test > output
done
```
输出文件`output`：
```bash
line 10
line 11
line 12
line 13
line 14
line 15
line 16
line 17
line 18
line 19
line 20
line 21
line 22
line 23
line 24
line 25
line 26
line 27
line 28
line 29
```
测试六：
```bash
$ ./selpg -s 2 -e 0 -l 5 test
invalid end page 0

Usage: selpg [-s start_page] [-e end_page] [ -f | -l lines_per_page ] [ -d dest ] [ in_filename ]

Options:
  -d, --d string
    	destination
  -e, --e int
    	end page (default -1)
  -f, --f
    	form-feed-delimited
  -l, --l int
    	lines per page (default 72)
  -s, --s int
    	start page (default -1)
```