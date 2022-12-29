// Get value from buttinon and set it to the modal-delete.html dialog form
// Make sure the buttons have id Example: <button id="trigger-modal-delete{{.ID}}" type="button"> 
var buttons = document.querySelectorAll("button[id^=trigger-modal-delete]");

buttons.forEach(btn => {
    btn.addEventListener('click', event => {
        var values = event.target.value.split(";")// "ID;Text"

        document.getElementById("ID").value = values[0];
        document.getElementById("modal-delete-label").innerHTML = values[0];
        document.getElementById("modal-delete-value").value = values[1];
    });
 });

// For Modal Update
 var editButtons = document.querySelectorAll("button[id^=trigger-modal-update]");

 editButtons.forEach(btn => {
    btn.addEventListener('click', event => {
        var values = event.target.value.split(";") //"ID;V0;V1"
        var id = values[0];
        var v0 = values[1];

        document.getElementById("ID").value = id;
        document.getElementById("modal-update-label").innerHTML = v0;
    });
 });

 // Page: /admin/users/list
 // Action: checkbox enables/disables delete-all-button
 var allCheckboxex = document.querySelectorAll(".form-check-input");
 var ckbCounter = 0;

 allCheckboxex.forEach(btn => {
    btn.addEventListener('click', event => {
        if (event.target.checked) {
            ckbCounter++;
        } else {
            ckbCounter--;
        }
            if (ckbCounter>0) {
                var deleteAllBotton = document.querySelector("button[id^=delete-all-button");
                deleteAllBotton.disabled = false
             } else {
                var deleteAllBotton = document.querySelector("button[id^=delete-all-button");
                deleteAllBotton.disabled = true
             }
        });
    });
    

 