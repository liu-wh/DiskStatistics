import { Row, Col, Result, Card } from 'antd';
import { Treemap } from '@ant-design/plots';


function byteConvert(bytes) {
    if (isNaN(bytes)) {
        return '';
    }
    let symbols = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    let exp = Math.floor(Math.log(bytes)/Math.log(2));
    if (exp < 1) {
        exp = 0;
    }
    let i = Math.floor(exp / 10);
    bytes = bytes / Math.pow(2, 10 * i);

    if (bytes.toString().length > bytes.toFixed(2).toString().length) {
        bytes = bytes.toFixed(2);
    }
    return bytes + ' ' + symbols[i];

}


function ScanResult(props) {

    const {data} = props



    const treeMapConfig = {
        data: data,
        colorField: 'name',
        legend: {
            position: 'bottom',
        },
        label: {
            style: {
                fill: 'black',
                fontSize: 14
            },
        },
        // use `drilldown: { enabled: true }` to
        // replace `interactions: [{ type: 'treemap-drill-down' }]`
        tooltip: {
            formatter: (v) => {
                const root = v.path[v.path.length - 1];
                return {
                    name: v.name,
                    value: `${byteConvert(v.value)} (${((v.value / root.value) * 100).toFixed(2)}%)`,
                };
            },
        },
        interactions: [
            {
                type: 'treemap-drill-down',
            },
        ],
        drilldown: {
            enabled: true,
            breadCrumb: {
                rootText: 'C:',
            },
        },
        animation: {},
    };

    return (
        <>
            <Treemap   {...treeMapConfig} />
        </>

    )

}

export default ScanResult