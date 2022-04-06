/*

*/

$(document).ready(function(){

    /*on() is used instead of click because click can be used only on static elements, and on() is to be used when you add
    elements dynamically*/
    $('[data-toggle="tooltip"]').tooltip();
  
      // Build out modal for deleting objects based on bootstrap attributes
    var addObjectModal = document.getElementById('addObjectModal')
    addObjectModal.addEventListener('show.bs.modal', function (event) {
      // Button that triggered the modal
      var button = event.relatedTarget
      // Extract info from data-bs-* attributes
      var objectType = button.getAttribute('data-bs-objecttype')
      //Setup form action based on Object Type
      if (objectType == "group") {
        $("addObjectForm").attr("action","/groups/add");
      } else if (objectType == "user") {
        $("#addObjectForm").attr("action","/users/add");
      } else if (objectType == "safe") {
        $("#addObjectForm").attr("action","/safes/add");
      } else {}
  
    });
  
    // Build out modal for deleting objects based on bootstrap attributes
    var deleteObjectModal = document.getElementById('deleteObjectModal')
    deleteObjectModal.addEventListener('show.bs.modal', function (event) {
      // Button that triggered the modal
      var button = event.relatedTarget
      // Extract info from data-bs-* attributes
      var objectId= button.getAttribute('data-bs-id')
      var objectType = button.getAttribute('data-bs-objecttype')
      var displayName = button.getAttribute('data-bs-displayname')
  
      // Update the modal's content.
      var modalBodyWarning = deleteObjectModal.querySelector('.modal-body #delete-warning')
      modalBodyWarning.textContent = 'Delete ' + objectType + ' "' + displayName + '"?'
      
      //Setup form action based on Object Type and ID
      if (objectType == "group") {
        $("#delObjectForm").attr("action","/groups/del/" + objectId);
      } else if (objectType == "user") {
        $("#delObjectForm").attr("action","/users/del/" + objectId);
      } else if (objectType == "safe") {
        $("#delObjectForm").attr("action","/safes/del/" + objectId);
      } else {}
  
    });
  
    // Build out modal for Updating objects based on bootstrap attributes
    var updateObjectModal = document.getElementById('updateObjectModal')
    updateObjectModal.addEventListener('show.bs.modal', function (event) {    
      // Button that triggered the modal
      var button = event.relatedTarget
      // Extract info from data-bs-* attributes
      var objectId= button.getAttribute('data-bs-id')
      var objectType = button.getAttribute('data-bs-objecttype')
      var displayName = button.getAttribute('data-bs-displayname')
  
      // Update the modal's content.
      var modalBodyWarning = updateObjectModal.querySelector('.modal-body #update-warning')
      modalBodyWarning.textContent = 'Update ' + objectType + ' "' + displayName + '"?'
      
      //Setup form action based on Object Type and ID
      if (objectType == "group") {
        $("#delObjectForm").attr("action","/groups/update/" + objectId);
      } else if (objectType == "user") {
        $("#delObjectForm").attr("action","/users/update/" + objectId);
      } else if (objectType == "safe") {
        $("#delObjectForm").attr("action","/safes/update/" + objectId);
      } else {}
  
    });
  
    // Build out modal for Inspecting objects based on bootstrap attributes
    var inspectObjectModal = document.getElementById('inspectObjectModal')
    inspectObjectModal.addEventListener('show.bs.modal', function (event) {
      // Button that triggered the modal
      var button = event.relatedTarget
      // Extract info from data-bs-* attributes
      var objectId= button.getAttribute('data-bs-id')
      var objectType = button.getAttribute('data-bs-objecttype')
      var displayName = button.getAttribute('data-bs-displayname')
      //var requestName = button.getAttribute('data-bs-requestname')
  
      // Setup form action based on Object Type and ID
      if (objectType == "group") {
        var urlPath = 'groups/' + displayName
      } else if (objectType == "user") {
        // User query based on Display Name
        var urlPath = 'users/' + displayName
      } else if (objectType == "Safe") {
        // Safe query based on Display Name
        var urlPath = 'safes/' + displayName
      } else {}
  
      const url = 'https://scimplistic.strlab.us/' + urlPath;
      
      getData(url);
      
    });
    
    $("#noti").click(
        function(){
            this.fadeOut();
        }  
    );
    if ($('#actlMsg').html()==''){
        $('.notification').addClass('hidden');
    } else {
      $('.notification').fadeOut(9000);
    }
    $('.btnMessage').click(function(){$('.notification').fadeOut()})
  
      
    $("#addUserBtn").on("click", function() {
          /*this.preventDefaults();
          var task_id = $("#task-id").val();
          $.ajax({
              url: "/tasks/" + task_id,
              type: "POST",
          data: {'title':'randome note', 'content':'this and that'}
          }).done(function(res, status) {
              console.log(status, res);
          var response = res
          $("#timeline").append(response)
          });*/
    });
  
    $('.toggle').click(function(){
        $(this).next().toggle();
      });
  
    async function getData(url) {
      // Default options are marked with *
      const response = await fetch(url)
      // Storing data in form of JSON
      var data = await response.json();
      console.log(data);
      if (response) {
          hideloader();
      }
    return response.json(); // parses JSON response into native JavaScript objects    
    }
  
  // Function to hide the loader
  function hideloader() {
    document.getElementById('loading').style.display = 'none';
  }
  
  });