/*

*/

$(document).ready(function(){

  /*on() is used instead of click because click can be used only on static elements, and on() is to be used when you add
  elements dynamically*/
  $('[data-toggle="tooltip"]').tooltip();

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

  $('.toggle').click(function(){
      $(this).next().toggle();
    });
});