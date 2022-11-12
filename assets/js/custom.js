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