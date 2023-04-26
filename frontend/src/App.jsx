import {useState} from 'react';
import {DiskTreeMapStatistics} from "../wailsjs/go/main/App";
import { Line } from '@ant-design/charts';
import ScanResult from "./Componets/index.jsx";






function App() {

    const [data, setData ] =  useState(false)
    const [lock, setLock] = useState(false)


    const onClick = ()=> {
        DiskTreeMapStatistics("root").then(
            (data) =>{
                if (lock) {
                    return
                }
                setLock(true)
                setData(data)
                setLock(false)
            }
        )
    }

    return (
        <div style={{height: 700, width: 1024, background: "white"}}>
            {
                data === false ?
                    <button onClick={onClick}>开始扫描</button>
                     : <ScanResult data={data}/>

            }
        </div>
    )
}

export default App
