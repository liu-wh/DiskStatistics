import {useState} from 'react';
import {DiskTreeMapStatistics} from "../wailsjs/go/main/App";
import { Button, Result, message } from 'antd';
import ScanResult from "./Componets/index.jsx";
import 'antd/dist/antd.css';






function App() {

    const [data, setData ] =  useState(false)
    const [lock, setLock] = useState(false)
    const [messageApi, contextHolder] = message.useMessage();


    const onClick = ()=> {
        messageApi.open({
            type: 'loading',
            content: '正在扫描中, 请稍候.....',
            duration: 0,
            key: "loading",
        })
        DiskTreeMapStatistics("root").then(
            (data) =>{
                if (lock) {
                    return
                }
                setLock(true)
                setData(data)
                setLock(false)
                messageApi.open({
                    key: "loading",
                    type: 'success',
                    content: '扫描成功!',
                    duration: 1,
                })
            }
        )
    }

    return (
        <div  >
            {contextHolder}
            {
                data === false ?
                    <Result
                        title="还没有进行扫描"
                        extra={
                            <Button type="primary" key="console" onClick={onClick}>
                                开始扫描
                            </Button>
                        }
                    />
                     : <ScanResult data={data}/>

            }
        </div>
    )
}

export default App
