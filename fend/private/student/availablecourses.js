fetch('/api/user/private/availablecourse').then((data)=>{
    return data.json();
  }).then((completedata)=>{
    console.log(completedata)
    let data1="";
      completedata.map((values)=>{
          data1+=`
            <div class="card">
              <div class="card-body">
              <h5 class="card-title" id="title">${values.title}</h5>
              <h5 class="card-title" id="teacher">${values.teacher}</h5>
              <p class="card-text" id="description">${values.description}</p>
              <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#myModal">Enroll</button>

            </div>
            </div>
          `
      });
      document.getElementById("card").innerHTML=data1;
  }).catch((err)=>{
    console.log(err)
  })