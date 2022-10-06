// an empty dropdown menu
const taskSelector = document.getElementById('taskSelector');

// empty div
const taskDisplayer = document.getElementById('taskDisplayer');


const taskFetcherBtn = document.getElementById('taskFetcherBtn');


const displayTask = (array) => {
    // clear the div
    taskDisplayer.innerHTML = '';
    // clear the taskDisplayer
    
    // create a new h1 element
    let h1 = createElement('p', array[0]);
    // create a new p element
    let p1 = createElement('p', array[1]);
    // create a new p element
    let answer = createElement('p', array[2]);
    // append the h1 element to the taskDisplayer
    taskDisplayer.appendChild(h1);
    // append the p element to the taskDisplayer
    taskDisplayer.appendChild(p1);
    // append the p element to the taskDisplayer
    taskDisplayer.appendChild(answer);
}



const fetchNewTask = async () => {
    const response = await fetch('/tasks');
    const data = await response.json();
    return data.tasks;
}

const addValuesToDropdown = (array) => {
    console.log(array);
    // get the first element of the array and add it to the dropdown menu as an option, the values is the entire array

    // add the first element of the array to the dropdown menu
    let option = createElement('option', array[0]);
    // set the value of the option to the entire array
    option.value = array;
    // append the option to the dropdown menu
    taskSelector.appendChild(option);
}

const getSelectedTask = () => {

    // get the value of the selected option
    let selected = taskSelector.value;
    //separate the string into an array
    let array = selected.split(',');

    
    displayTask(array);
}

const createElement = (element, text) => {
    let newElement = document.createElement(element);
    newElement.innerHTML = text;
    return newElement;
}

taskFetcherBtn.addEventListener('click', async () => {
    const tasks = await fetchNewTask();

    // display the new task
    displayTask(tasks);

        addValuesToDropdown(tasks);
    
}
)

taskSelector.addEventListener('change', () => {
    getSelectedTask();
}
)

//on page load, fetch one task and display it
window.onload = async () => {
    const tasks = await fetchNewTask();

    addValuesToDropdown(tasks);

    displayTask(tasks);
}