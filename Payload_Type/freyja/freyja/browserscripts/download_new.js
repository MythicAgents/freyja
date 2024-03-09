function(task, responses){

    if(task.status.includes("error")){
        const combined = responses.reduce( (prev, cur) => {
            return prev + cur;
        }, "");
        return {'plaintext': combined};
    }else if(task.completed){
        if(responses.length > 0){
            try{
                let data = JSON.parse(responses[0]);
                return { "media": [{
                        "filename": `${task.display_params}`,
                        "agent_file_id": data["file_id"],
                    }]};
            }catch(error){
                const combined = responses.reduce( (prev, cur) => {
                    return prev + cur;
                }, "");
                return {'plaintext': combined};
            }

        }else{
            return {"plaintext": "No data to display..."}
        }

    }else{
        if(responses.length > 0){
            const task_data = JSON.parse(responses[0]);
            return {"plaintext": "Downloading a file with " + task_data["total_chunks"] + " total chunks..."};
        }
        return {"plaintext": "No data yet..."}
    }
}