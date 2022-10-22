

fetch('/api/user/private/getyourcourse').then((data)=>{
  return data.json();
}).then((completedata)=>{
  console.log(completedata)
  let data1="";
    completedata.map((values)=>{
        data1+=`
          <div class="card">
            <div class="card-body">
            <h5 class="card-title" id="title">${values.title}</h5>
            <p class="card-text" id="description">${values.description}</p>
            <a href="#" id="${values.id}" class="btn btn-primary" >Resources</a>
            <a   id="${values.id}" onclick="deletecard(this.id)" class="btn btn-primary">Delete</a>
          </div>
          </div>
        `
    });
    document.getElementById("card").innerHTML=data1;
}).catch((err)=>{
  console.log(err)
})

function deletecard(index){
  console.log("I am deleting course card of id ",index);
  const id=Number(index);
  console.log(typeof(id))
  
  const postMethod = {
    method: 'POST', // Method itself
    headers: {
     'Content-type': 'application/json; charset=UTF-8' // Indicates the content 
    },
    body: JSON.stringify({
      id : id,
    }
    ) // We send data in JSON format
   }
   
   // make the HTTP put request using fetch api
   fetch(`/api/user/private/deletecourse`, postMethod)
   .then(response => response.json())
   .then(data => console.log(data)) // Manipulate the data retrieved back, if we want to do something with it
   .catch(err => console.log(err)) // Do something with the error
   location.reload();
}