fetch("/stats/json")
    .then((resp) => { return resp.json(); })
    .then((dataJSON) => {
        console.debug(dataJSON);

        document.getElementById("startTime").innerHTML = dataJSON.start_time;
        document.getElementById("cpuCount").innerHTML = dataJSON.num_cpu;
        document.getElementById("goroutineCount").innerHTML = dataJSON.num_goroutines;
        document.getElementById("memUsage").innerHTML = dataJSON.total_memory_alloc;
        document.getElementById("nodeCount").innerHTML = dataJSON.node_count;
        document.getElementById("edgeCount").innerHTML = dataJSON.edge_count;
    })