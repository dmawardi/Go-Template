// Select all the rows in the table
function selectAllRows(checkbox) {
  // Grab all the checkboxes
  var rowCheckboxes = document.querySelectorAll(".select-row");
  //   Iterate through the checkboxes and set checked to the value of the select-all checkbox
  for (var i = 0; i < rowCheckboxes.length; i++) {
    rowCheckboxes[i].checked = checkbox.checked;
  }
}

// Update the select-all checkbox status whenever a row checkbox is clicked
function updateSelectAll(checkbox) {
  // Grab the select-all checkbox
  var selectAllCheckbox = document.getElementById("select-all");
  //   If the row checkbox is not checked, then the select-all checkbox should not be checked
  if (!checkbox.checked) {
    selectAllCheckbox.checked = false;
  } else {
    //  If the row checkbox is checked, then check if all the row checkboxes are checked
    var allRowCheckboxesChecked = true;
    // Grab all the checkboxes
    var rowCheckboxes = document.querySelectorAll(".select-row");
    // Iterate through the checkboxes and set allRowCheckboxesChecked to false if any of the row checkboxes are not checked
    for (var i = 0; i < rowCheckboxes.length; i++) {
      if (!rowCheckboxes[i].checked) {
        allRowCheckboxesChecked = false;
        break;
      }
    }
    // Set the select-all checkbox to the value of allRowCheckboxesChecked
    selectAllCheckbox.checked = allRowCheckboxesChecked;
  }
}
