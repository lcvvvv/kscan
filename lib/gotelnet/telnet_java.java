using System;
using System.Collections.Generic;
using System.Text;
using System.Net;
using System.Net.Sockets;
using System.Collections;

namespace ConsoleApplication1
{
    public class Program
    {
        #region 一些telnet的数据定义,先没看懂没关系
        /// <summary>
        /// 标志符,代表是一个TELNET 指令
        /// </summary>
        readonly Char IAC = Convert.ToChar(255);
        /// <summary>
        /// 表示一方要求另一方使用，或者确认你希望另一方使用指定的选项。
        /// </summary>
        readonly Char DO = Convert.ToChar(253);
        /// <summary>
        /// 表示一方要求另一方停止使用，或者确认你不再希望另一方使用指定的选项。
        /// </summary>
        readonly Char DONT = Convert.ToChar(254);
        /// <summary>
        /// 表示希望开始使用或者确认所使用的是指定的选项。
        /// </summary>
        readonly Char WILL = Convert.ToChar(251);
        /// <summary>
        /// 表示拒绝使用或者继续使用指定的选项。
        /// </summary>
        readonly Char WONT = Convert.ToChar(252);

        /// <summary>
        /// 表示后面所跟的是对需要的选项的子谈判
        /// </summary>
        readonly Char SB = Convert.ToChar(250);

        /// <summary>
        /// 子谈判参数的结束
        /// </summary>
        readonly Char SE = Convert.ToChar(240);

        const Char IS = '0';

        const Char SEND = '1';

        const Char INFO = '2';

        const Char VAR = '0';

        const Char VALUE = '1';

        const Char ESC = '2';

        const Char USERVAR = '3';

        /// <summary>
        /// 流
        /// </summary>
        byte[] m_byBuff = new byte[100000];

        /// <summary>
        /// 收到的控制信息
        /// </summary>
        private ArrayList m_ListOptions = new ArrayList();

        /// <summary>
        /// 存储准备发送的信息
        /// </summary>
        string m_strResp;

        /// <summary>
        /// 一个Socket套接字
        /// </summary>
        private Socket s;
        #endregion


        /// <summary>
        /// 主函数
        /// </summary>
        /// <param name="args"></param>
        static void Main(string[] args)
        {
            //实例化这个对象
            Program p = new Program();
            //启动socket进行telnet 链接
            p.doSocket();


        }

        /// <summary>
        /// 启动socket 进行telnet操作
        /// </summary>
        private void doSocket()
        {
            //获得链接的地址,可以是网址也可以是IP
            Console.WriteLine("Server Address:");
            //解析输入,如果是一个网址,则解析成ip
            IPAddress import = GetIP(Console.ReadLine());
            //获得端口号
            Console.WriteLine("Server Port:");
            int port = int.Parse(Console.ReadLine());

            //建立一个socket对象,使用IPV4,使用流进行连接,使用tcp/ip 协议
            s = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

            //获得一个链接地址对象(由IP地址和端口号构成)
            IPEndPoint address = new IPEndPoint(import, port);

            /*
             * 说明此socket不是处于阻止模式
             *
             * msdn 对阻止模式的解释:
             * ============================================================
             * 如果当前处于阻止模式，并且进行了一个并不立即完成的方法调用，
             * 则应用程序将阻止执行，直到请求的操作完成后才解除阻止。
             * 如果希望在请求的操作尚未完成的情况下也可以继续执行，
             * 请将 Blocking 属性更改为 false。Blocking 属性对异步方法无效。
             * 如果当前正在异步发送和接收数据，并希望阻止执行，
             * 请使用 ManualResetEvent 类。
             * ============================================================
            */
            s.Blocking = false;

            /*
             * 开始一个对远程主机连接的异步请求,
             * 因为Telnet 使用的是TCP 链接,是面向连接的,
             * 所以此处BeginConnect 会启动一个异步请求,
             * 请求获得与 给的address 的连接
             *
             * 此方法的第二个函数是一个类型为AsyncCallback 的委托
             *
             * 这个AsyncCallback msdn给出的定义如下
             * ===================================================================
             * 使用 AsyncCallback 委托在一个单独的线程中处理异步操作的结果。A
             * syncCallback 委托表示在异步操作完成时调用的回调方法。
             * 回调方法采用 IAsyncResult 参数，该参数随后可用来获取异步操作的结果。
             * ===================================================================
             * 这个方法里的委托实际上就是 当异步请求有回应了之后,执行委托的方法.
             * 委托里的参数,实际上就是BeginConnect的第三个参数,
             * 此处为socket 本身
             *
             * 我比较懒,写了一个匿名委托,实际上跟AsyncCallback 效果一个样.
             *
             */
            s.BeginConnect(
                address,
                delegate(IAsyncResult ar)
                /*
                 * 此处为一个匿名委托,
                 * 实际上等于
                 * 建立一个AsyncCallback对象,指定后在此引用一个道理
                 *
                 * ok这里的意义是,
                 * 当远程主机连接的异步请求有响应的时候,执行以下语句
                 */
                {
                    try
                    {
                        //获得传入的对象 (此处对象是BeginConnect 的第三个参数)
                        Socket sock1 = (Socket)ar.AsyncState;

                        /*
                         * 如果 Socket 在最近操作时连接到远程资源，则为 true；否则为 false。
                         *
                         * 以下是MSDN 对Connected属性的备注信息
                         * =========================================================================
                         * Connected 属性获取截止到最后的 I/O 操作时 Socket 的连接状态。
                         * 当它返回 false 时，表明 Socket 要么从未连接，要么已断开连接。
                         *
                         * Connected 属性的值反映最近操作时的连接状态。如果您需要确定连接的当前状态，
                         * 请进行非阻止、零字节的 Send 调用。
                         * 如果该调用成功返回或引发 WAEWOULDBLOCK 错误代码 (10035)，
                         * 则该套接字仍然处于连接状态；否则，该套接字不再处于连接状态。
                         * =========================================================================
                        */
                        if (sock1.Connected)
                        {

                            AsyncCallback recieveData = new AsyncCallback(OnRecievedData);
                            /*
                             * 此处没再用匿名委托的原因是,
                             * 一个匿名委托嵌套一个匿名委托,我自己思路跟不上来了...
                             *
                             * ok,这里是当Connected  为true时,
                             * 使用BeginReceive 方法
                             * 开始接收信息到m_byBuff(我们在类中定义的私有属性)
                             *
                             * 以下是MSDN 对BeginReceive 的一些说明
                             * =========================================================================
                             * 异步 BeginReceive 操作必须通过调用 EndReceive 方法来完成。
                             * 通常，该方法由 callback 委托调用。此方法在操作完成前不会进入阻止状态。
                             * 若要一直阻塞到操作完成时为止，请使用 Receive 方法重载中的一个。
                             * 若要取消挂起的 BeginReceive，请调用 Close 方法。
                             * ==========================================================================
                             *
                             * 当接收完成之后,他们就会调用OnRecievedData方法
                             * 我在recieveData所委托的方法OnRecievedData 中调用了sock.EndReceive(ar);
                             */
                            sock1.BeginReceive(m_byBuff, 0, m_byBuff.Length, SocketFlags.None, recieveData, sock1);
                        }
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine("初始化接收信息出错:" + ex.Message);
                    }
                },
                s);

            //此处是为了发送指令而不停的循环
            while (true)
            {
                //发送读出的数据
                DispatchMessage(Console.ReadLine());
                //因为每发送一行都没有发送回车,故在此处补上
                DispatchMessage("\r\n");
            }

        }

        /// <summary>
        /// 当接收完成后,执行的方法(供委托使用)
        /// </summary>
        /// <param name="ar"></param>
        private void OnRecievedData(IAsyncResult ar)
        {
            //从参数中获得给的socket 对象
            Socket sock = (Socket)ar.AsyncState;
            /*
             * EndReceive方法为结束挂起的异步读取
             * (貌似是在之前的beginReceive收到数据之后,
             * socket只是"挂起",并未结束)
             * 之后返回总共接收到的字流量
             *
             * 以下是MSDN给出的EndReceive 的注意事项
             * =========================================================================================
             * EndReceive 方法完成在 BeginReceive 方法中启动的异步读取操作。
             *
             * 在调用 BeginReceive 之前，需创建一个实现 AsyncCallback 委托的回调方法。
             * 该回调方法在单独的线程中执行并在 BeginReceive 返回后由系统调用。
             * 回调方法必须接受 BeginReceive 方法所返回的 IAsyncResult 作为参数。
             *
             * 在回调方法中，调用 IAsyncResult 的 AsyncState 方法以获取传递给 BeginReceive 方法的状态对象。
             * 从该状态对象提取接收 Socket。在获取 Socket 之后，可以调用 EndReceive 方法以成功完成读取操作，
             * 并返回已读取的字节数。
             *
             * EndReceive 方法将一直阻止到有数据可用为止。
             * 如果您使用的是无连接协议，则 EndReceive 将读取传入网络缓冲区中第一个排队的可用数据报。
             * 如果您使用的是面向连接的协议，则 EndReceive 方法将读取所有可用的数据，
             * 直到达到 BeginReceive 方法的 size 参数所指定的字节数为止。
             * 如果远程主机使用 Shutdown 方法关闭了 Socket 连接，并且所有可用数据均已收到，
             * 则 EndReceive 方法将立即完成并返回零字节。
             *
             * 若要获取接收到的数据，请调用 IAsyncResult 的 AsyncState 方法，
             * 然后提取所产生的状态对象中包含的缓冲区。
             *
             * 若要取消挂起的 BeginReceive，请调用 Close 方法。
             * =========================================================================================
             */
            int nBytesRec = sock.EndReceive(ar);
            //如果有接收到数据的话
            if (nBytesRec > 0)
            {
                //将接收到的数据转个码,顺便转成string型
                string sRecieved = Encoding.GetEncoding("utf-8").GetString(m_byBuff, 0, nBytesRec);

                //声明一个字符串,用来存储解析过的字符串
                string m_strLine = "";
                //遍历Socket接收到的字符

                /*
                 * 此循环用来调整linux 和 windows在换行上标记的区别
                 * 最后将调整好的字符赋予给 m_strLine
                */
                for (int i = 0; i < nBytesRec; i++)
                {
                    Char ch = Convert.ToChar(m_byBuff[i]);
                    switch (ch)
                    {
                        case '\r':
                            m_strLine += Convert.ToString("\r\n");
                            break;
                        case '\n':
                            break;
                        default:
                            m_strLine += Convert.ToString(ch);
                            break;
                    }
                }

                try
                {
                    //获得转义后的字符串的长度
                    int strLinelen = m_strLine.Length;
                    //如果长度为零
                    if (strLinelen == 0)
                    {
                        //则返回"\r\n" 即回车换行
                        m_strLine = Convert.ToString("\r\n");
                    }

                    //建立一个流,把接收的信息(转换后的)存进 mToProcess 中
                    Byte[] mToProcess = new Byte[strLinelen];
                    for (int i = 0; i < strLinelen; i++)
                        mToProcess[i] = Convert.ToByte(m_strLine[i]);

                    // Process the incoming data
                    //对接收的信息进行处理,包括对传输过来的信息的参数的存取和
                    string mOutText = ProcessOptions(mToProcess);
                    //解析命令后返回 显示信息(即除掉了控制信息)
                    if (mOutText != "")
                        Console.Write(mOutText);


                    // Respond to any incoming commands
                    //接收完数据,处理完字符串数据等一系列事物之后,开始回发数据
                    RespondToOptions();
                }
                catch (Exception ex)
                {
                    throw new Exception("接收数据的时候出错了! " + ex.Message);
                }
            }
            else// 如果没有接收到任何数据的话
            {
                // 输出   关闭连接
                Console.WriteLine("Disconnected", sock.RemoteEndPoint);
                // 关闭socket
                sock.Shutdown(SocketShutdown.Both);
                sock.Close();
                Console.Write("Game Over");
                Console.ReadLine();
            }
        }

        /// <summary>
        ///  发送数据的函数
        /// </summary>

        private void RespondToOptions()
        {
            try
            {
                //声明一个字符串,来存储 接收到的参数
                string strOption;
                /*
                 * 此处的控制信息参数,是之前接受到信息之后保存的
                 * 例如 255   253   23   等等
                 * 具体参数的含义需要去查telnet 协议
                 */
                for (int i = 0; i < m_ListOptions.Count; i++)
                {
                    //获得一个控制信息参数
                    strOption = (string)m_ListOptions[i];
                    //根据这个参数,进行处理
                    ArrangeReply(strOption);
                }
                DispatchMessage(m_strResp);
                m_strResp = "";
                m_ListOptions.Clear();
            }
            catch (Exception ers)
            {
                Console.WriteLine("错错了,在回发数据的时候 " + ers.Message);
            }
        }

        /// <summary>
        /// 解析接收的数据,生成最终用户看到的有效文字,同时将附带的参数存储起来
        /// </summary>
        /// <param name="m_strLineToProcess">收到的处理后的数据</param>
        /// <returns></returns>
        private string ProcessOptions(byte[] m_strLineToProcess)
        {

            string m_DISPLAYTEXT = "";
            string m_strTemp = "";
            string m_strOption = "";
            string m_strNormalText = "";
            bool bScanDone = false;
            int ndx = 0;
            int ldx = 0;
            char ch;
            try
            {
                //把数据从byte[] 转化成string
                for (int i = 0; i < m_strLineToProcess.Length; i++)
                {
                    Char ss = Convert.ToChar(m_strLineToProcess[i]);
                    m_strTemp = m_strTemp + Convert.ToString(ss);
                }

                //此处意义为,当没描完数据前,执行扫描
                while (bScanDone != true)
                {
                    //获得长度
                    int lensmk = m_strTemp.Length;
                    //之后开始分析指令,因为每条指令为255 开头,故可以用此来区分出每条指令
                    ndx = m_strTemp.IndexOf(Convert.ToString(IAC));

                    //此处为出错判断,本无其他含义
                    if (ndx > lensmk)
                        ndx = m_strTemp.Length;

                    //此处为,如果搜寻到IAC标记的telnet 指令,则执行以下步骤
                    if (ndx != -1)
                    {
                        #region 如果存在IAC标志位
                        // 将 标志位IAC 的字符 赋值给最终显示文字
                        m_DISPLAYTEXT += m_strTemp.Substring(0, ndx);
                        // 此处获得命令码
                        ch = m_strTemp[ndx + 1];

                        //如果命令码是253(DO) 254(DONT)  521(WILL) 252(WONT) 的情况下
                        if (ch == DO || ch == DONT || ch == WILL || ch == WONT)
                        {
                            //将以IAC 开头3个字符组成的整个命令存储起来
                            m_strOption = m_strTemp.Substring(ndx, 3);
                            m_ListOptions.Add(m_strOption);

                            // 将 标志位IAC 的字符 赋值给最终显示文字
                            m_DISPLAYTEXT += m_strTemp.Substring(0, ndx);

                            //将处理过的字符串删去
                            string txt = m_strTemp.Substring(ndx + 3);
                            m_strTemp = txt;
                        }
                        //如果IAC后面又跟了个IAC (255)
                        else if (ch == IAC)
                        {
                            //则显示从输入的字符串头开始,到之前的IAC 结束
                            m_DISPLAYTEXT = m_strTemp.Substring(0, ndx);
                            //之后将处理过的字符串排除出去
                            m_strTemp = m_strTemp.Substring(ndx + 1);
                        }
                        //如果IAC后面跟的是SB(250)
                        else if (ch == SB)
                        {
                            m_DISPLAYTEXT = m_strTemp.Substring(0, ndx);
                            ldx = m_strTemp.IndexOf(Convert.ToString(SE));
                            m_strOption = m_strTemp.Substring(ndx, ldx);
                            m_ListOptions.Add(m_strOption);
                            m_strTemp = m_strTemp.Substring(ldx);
                        }

                        #endregion
                    }
                    //若字符串里已经没有IAC标志位了
                    else
                    {
                        //显示信息累加上m_strTemp存储的字段
                        m_DISPLAYTEXT = m_DISPLAYTEXT + m_strTemp;
                        bScanDone = true;
                    }
                }
                //输出人看到的信息
                m_strNormalText = m_DISPLAYTEXT;
            }
            catch (Exception eP)
            {
                throw new Exception("解析传入的字符串错误:" + eP.Message);
            }
            return m_strNormalText;

        }

        /// <summary>
        /// 获得IP地址
        /// </summary>
        /// <param name="import"></param>
        /// <returns></returns>
        private static IPAddress GetIP(string import)
        {
            IPHostEntry IPHost = Dns.GetHostEntry(import);
            return IPHost.AddressList[0];
        }




        #region magic Function

        //解析传过来的参数,生成回发的数据到m_strResp
        private void ArrangeReply(string strOption)
        {
            try
            {

                Char Verb;
                Char Option;
                Char Modifier;
                Char ch;
                bool bDefined = false;
                //排错选项,无啥意义
                if (strOption.Length < 3) return;

                //获得命令码
                Verb = strOption[1];
                //获得选项码
                Option = strOption[2];

                //如果选项码为 回显(1) 或者是抑制继续进行(3)
                if (Option == 1 || Option == 3)
                {
                    bDefined = true;
                }
                // 设置回发消息,首先为标志位255
                m_strResp += IAC;

                //如果选项码为 回显(1) 或者是抑制继续进行(3) ==true
                if (bDefined == true)
                {
                    #region 继续判断
                    //如果命令码为253 (DO)
                    if (Verb == DO)
                    {
                        //我设置我应答的命令码为 251(WILL) 即为支持 回显或抑制继续进行
                        ch = WILL;
                        m_strResp += ch;
                        m_strResp += Option;

                    }
                    //如果命令码为 254(DONT)
                    if (Verb == DONT)
                    {
                        //我设置我应答的命令码为 252(WONT) 即为我也会"拒绝启动" 回显或抑制继续进行
                        ch = WONT;
                        m_strResp += ch;
                        m_strResp += Option;

                    }
                    //如果命令码为251(WILL)
                    if (Verb == WILL)
                    {
                        //我设置我应答的命令码为 253(DO) 即为我认可你使用回显或抑制继续进行
                        ch = DO;
                        m_strResp += ch;
                        m_strResp += Option;
                        //break;
                    }
                    //如果接受到的命令码为251(WONT)
                    if (Verb == WONT)
                    {
                        //应答  我也拒绝选项请求回显或抑制继续进行
                        ch = DONT;
                        m_strResp += ch;
                        m_strResp += Option;
                        //    break;
                    }
                    //如果接受到250(sb,标志子选项开始)
                    if (Verb == SB)
                    {
                        /*
                         * 因为启动了子标志位,命令长度扩展到了4字节,
                         * 取最后一个标志字节为选项码
                         * 如果这个选项码字节为1(send)
                         * 则回发为 250(SB子选项开始) + 获取的第二个字节 + 0(is) + 255(标志位IAC) + 240(SE子选项结束)
                        */
                        Modifier = strOption[3];
                        if (Modifier == SEND)
                        {
                            ch = SB;
                            m_strResp += ch;
                            m_strResp += Option;
                            m_strResp += IS;
                            m_strResp += IAC;
                            m_strResp += SE;
                        }
                    }
                    #endregion
                }
                else //如果选项码不是1 或者3
                {
                    #region 底下一系列代表,无论你发那种请求,我都不干
                    if (Verb == DO)
                    {
                        ch = WONT;
                        m_strResp += ch;
                        m_strResp += Option;
                    }
                    if (Verb == DONT)
                    {
                        ch = WONT;
                        m_strResp += ch;
                        m_strResp += Option;
                    }
                    if (Verb == WILL)
                    {
                        ch = DONT;
                        m_strResp += ch;
                        m_strResp += Option;
                    }
                    if (Verb == WONT)
                    {
                        ch = DONT;
                        m_strResp += ch;
                        m_strResp += Option;
                    }
                    #endregion
                }
            }
            catch (Exception eeeee)
            {

                throw new Exception("解析参数时出错:" + eeeee.Message);
            }

        }

        /// <summary>
        /// 将信息转化成charp[] 流的形式,使用socket 进行发出
        /// 发出结束之后,使用一个匿名委托,进行接收,
        /// 之后这个委托里,又有个委托,意思是接受完了之后执行OnRecieveData 方法
        ///
        /// </summary>
        /// <param name="strText"></param>
        void DispatchMessage(string strText)
        {
            try
            {
                //申请一个与字符串相当长度的char流
                Byte[] smk = new Byte[strText.Length];
                for (int i = 0; i < strText.Length; i++)
                {
                    //解析字符串,将其存储到char流中去
                    Byte ss = Convert.ToByte(strText[i]);
                    smk[i] = ss;
                }

                //发送char流,之后发送完毕后执行委托中的方法(此处为匿名委托)
                /*MSDN 对BeginSend 的解释
                 * =======================================================================================================
                 * BeginSend 方法可对在 Connect、BeginConnect、Accept 或 BeginAccept 方法中建立的远程主机启动异步发送操作。
                 * 如果没有首先调用 Accept、BeginAccept、Connect 或 BeginConnect，则 BeginSend 将会引发异常。
                 * 调用 BeginSend 方法将使您能够在单独的执行线程中发送数据。
                 * 您可以创建一个实现 AsyncCallback 委托的回调方法并将它的名称传递给 BeginSend 方法。
                 * 为此，您的 state 参数至少必须包含用于通信的已连接或默认 Socket。
                 * 如果回调需要更多信息，则可以创建一个小型类或结构，用于保存 Socket 和其他所需的信息。
                 * 通过 state 参数将此类的一个实例传递给 BeginSend 方法。
                 * 回调方法应调用 EndSend 方法。
                 * 当应用程序调用 BeginSend 时，系统将使用一个单独的线程来执行指定的回调方法，
                 * 并阻止 EndSend，直到 Socket 发送了请求的字节数或引发了异常为止。
                 * 如果希望在调用 BeginSend 方法之后使原始线程阻止，请使用 WaitHandle.WaitOne 方法。
                 * 当需要原始线程继续执行时，请在回调方法中调用 T:System.Threading.ManualResetEvent 的 Set 方法。
                 * 有关编写回调方法的其他信息，请参见 Callback 示例。
                 * =======================================================================================================
                 */
                IAsyncResult ar2 = s.BeginSend(smk, 0, smk.Length, SocketFlags.None, delegate(IAsyncResult ar)
                {
                    //当执行完"发送数据" 这个动作后
                    // 获取Socket对象,对象从beginsend 中的最后个参数上获得
                    Socket sock1 = (Socket)ar.AsyncState;

                    if (sock1.Connected)//如果连接还是有效
                    {
                        //这里建立一个委托
                        AsyncCallback recieveData = new AsyncCallback(OnRecievedData);

                        /*
                         * 此处为:开始接受数据(在发送完毕之后-->出自于上面的匿名委托),
                         * 当接收完信息之后,执行OnrecieveData方法(由委托传进去),
                         * 注意,是异步调用
                         */
                        sock1.BeginReceive(m_byBuff, 0, m_byBuff.Length, SocketFlags.None, recieveData, sock1);
                    }
                }, s);

                /*
                 * 结束 异步发送
                 * EndSend 完成在 BeginSend 中启动的异步发送操作。
                 * 在调用 BeginSend 之前，需创建一个实现 AsyncCallback 委托的回调方法。
                 * 该回调方法在单独的线程中执行并在 BeginSend 返回后由系统调用。
                 * 回调方法必须接受 BeginSend 方法所返回的 IAsyncResult 作为参数。
                 *
                 * 在回调方法中，调用 IAsyncResult 参数的 AsyncState 方法可以获取发送 Socket。
                 * 在获取 Socket 之后，则可以调用 EndSend 方法以成功完成发送操作，并返回发送的字节数。
                 */
                s.EndSend(ar2);
            }
            catch (Exception ers)
            {
                Console.WriteLine("出错了,在回发数据的时候:" + ers.Message);
            }
        }
        #endregion
    }
}