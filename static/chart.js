let newData = {
    in_use: [],
};

const formatToMB = (value) => parseFloat(value / 1024 / 1024).toFixed(1);
const eventSource = new EventSource("/events");
const chart = new ApexCharts(document.getElementById("line-chart"), {
    tooltip: {
        enabled: true,
        x: { show: true },
        y: { show: true },
    },
    grid: {
        show: true,
        strokeDashArray: 4,
        padding: {
            left: 2,
            right: 2,
            top: -26,
        },
    },
    series: [
        {
            name: "In use",
            data: newData.in_use,
        },
    ],
    noData: { text: "Loading ..." },
    chart: {
        height: "100%",
        maxWidth: "100%",
        type: "area",
        animations: {
            enabled: true,
            easing: "linear",
            dynamicAnimation: { speed: 1 },
        },
        toolbar: { show: false },
    },
    fill: {
        // type: 'solid',
    },
    dataLabels: {
        enabled: false,
        type: "timestamp",
    },
    stroke: {
        width: 3,
        curve: "straight",
    },
    xaxis: {
        labels: { show: false },
        axisBorder: { show: false },
        axisTicks: { show: false },
    },
});
chart.render();

eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    const in_use = formatToMB(data.in_use);

    document.getElementById("total").innerHTML =
        formatToMB(data.total) + "MB";
    document.getElementById("in_use").innerHTML = in_use + "MB";
    document.getElementById("available").innerHTML =
        formatToMB(data.available) + "MB";

    let t = new Date().toLocaleTimeString();

    newData.in_use.push({ x: t, y: Number(in_use) });

    if (newData.in_use.length > 60) {
        newData.in_use.shift();
    }

    chart.updateSeries([{ data: newData.in_use }]);
};