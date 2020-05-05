var container = document.getElementById('graph');

var options = {
    height: '100%',
    width: '100%',
    nodes: {
        shape: 'circle',
        scaling: { min: 10, max: 20 },
        chosen: {
            node: (values, id, selected, hovering) => {
                values.color = "#ffe6e6";
                values.shadow = true;
            }
        },
    },
    edges: {
        chosen: {
            edge: (values, id, selected, hovering) => {
                values.color = "red";
            },
        },
    },
    physics: {
        solver: "forceAtlas2Based",
    },
};

var network = new vis.Network(container, [], options);

function convertJSON(data) {
    var options = {};
    var nodes = new vis.DataSet(options);
    var edges = new vis.DataSet(options);

    if (data.nodes == undefined) {
        data.nodes = [];
    }

    data.nodes.forEach(element => {
        var node = {
            id: element.uid,
            label: element.label,
            group: element.label,
            properties: {},
        }

        for (var key in element.properties) {
            node.properties[key] = atob(element.properties[key]);
        }

        console.debug(element);
        console.debug(node);
        nodes.add(node);
    });

    if (data.edges == undefined) {
        data.edges = [];
    }

    data.edges.forEach(element => {
        var edge = {
            id: element.uid,
            from: element.source_uid,
            label: element.label,
            to: element.target_uid,
            group: element.label,
            properties: {},
            arrows: "to",
        }

        for (var key in element.properties) {
            edge.properties[key] = atob(element.properties[key]);
        }

        console.debug(element);
        console.debug(edge);
        edges.add(edge);
    });

    return { "nodes": nodes, "edges": edges };
}

function GetQuery() {
    var queryTextArea = document.getElementById("cypher")
    return queryTextArea.name + "=" + queryTextArea.value;
}

function FetchFullGraphData() {
    fetch("/assets/json", { method: "GET", headers: { "Content-Type": "application/json" } })
        .then((resp) => { return resp.json() })
        .then((data) => { return convertJSON(data) })
        .then((store) => { network.setData(store) })
}

function FetchGraphData() {
    fetch("/assets/json", { method: "POST", headers: { "Content-Type": "application/x-www-form-urlencoded" }, body: GetQuery() })
        .then((resp) => {
            if (resp.ok) {
                return resp.json();
            } else {
                var notify = document.getElementById("query-toast");
                resp.text().then((text) => {
                    notify.MaterialSnackbar.showSnackbar(
                        {
                            message: text
                        }
                    )
                })

                return {};
            }
        })
        .then((dataJSON) => { return convertJSON(dataJSON); })
        .then((store) => {
            console.debug(store);
            network.setData(store);
        })
        .catch((error) => {
            var notify = document.getElementById("query-toast");
            notify.MaterialSnackbar.showSnackbar(
                {
                    message: error
                }
            )
        })
}
