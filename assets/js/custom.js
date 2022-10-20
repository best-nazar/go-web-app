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