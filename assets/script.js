const btn = document.querySelector('button');
const taskDiv = document.getElementById('task');

btn.addEventListener('click', async () => {

    //clear the task div
    taskDiv.innerHTML = '';
    
    //create a new child element for the task div
    const taskInput = document.createElement('p');
    const taskOutput = document.createElement('p');
    const taskResponse = document.createElement('p');
    taskDiv.appendChild(taskInput);
    taskDiv.appendChild(taskOutput);
    taskDiv.appendChild(taskResponse);


    await fetch('/tasks')
        .then(response => response.json())
        .then(data => {
            taskInput.innerHTML = data.tasks[0];
            taskOutput.innerHTML = data.tasks[1];
            taskResponse.innerHTML = data.tasks[2];
        }
    );
    

    });

