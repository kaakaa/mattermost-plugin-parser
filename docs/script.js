
var tabledata;
fetch("data.json")
    .then(res => res.json())
    .then(d => {
        var table = new Tabulator("#usage-table", {
            data:d,           //load row data from array
            layout:"fitColumns",      //fit columns to width of table
            responsiveLayout:"hide",  //hide columns that dont fit on the table
            tooltips:true,            //show tool tips on cells
            addRowPos:"top",          //when adding a new row, add it to the top of the table
            history:true,             //allow undo and redo actions on the table
            pagination:"local",       //paginate the data
            paginationSize: 100,         //allow 50 rows per page of data
            movableColumns:true,      //allow column order to be changed
            resizableRows:true,       //allow row order to be changed
            initialSort:[             //set the initial sort order of the data
                {column:"name", dir:"asc"},
            ],
            columns:[                 //define the table columns
                {title:"Name", field:"url", headerFilter: true, formatter: "link", formatterParams: {
                    labelField: "name",
                    target: "_blank",
                }},
                {title:"Refs", field:"refs"},
                {title:"Location", field:"loc", formatter:"html"},
                {title:"Type", field:"type"},
                {title:"API", field:"api", headerFilter: true},
                // {title:"CommitID", field:"commit_id"},
                {title:"Collected At", field: "created_at", formatter: "datetime", formatterParams: {
                    inputFormat: "YYYY-MM-DDTHH:mm:ss.SSSZ",
                    outputFormat: "YYYY-MM-DDTHH:mm:ss",
                    invalidPlaceholder: "(unknown)",
                }}
            ],
            footerElement: "<p color='white'>Develop at <a href='https://github.com/kaakaa/mattermost-plugin-parser'>kaakaa/mattermost-plugin-parser</a></p>",
        });        
    });
