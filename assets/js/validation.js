$(document).ready(function() { 
 
    $('#Singin').click(function() {  
 
        $(".error").hide();
        var hasError = false;
        var emailReg = /^([\w-\.]+@([\w-]+\.)+[\w-]{2,4})?$/;
        var passreg=/^(?=.*\d)(?=.*[A-Z])(?=.*[a-z]).{8,}$/;
        var passw = $("#Password").val();
        var emailaddressVal = $("#Emailid").val();
        if(emailaddressVal == '') {
            $("#Emailid").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Please Enter a your email address.</span>');
            hasError = true;
        }
 
        else if(!emailReg.test(emailaddressVal)) {
            $("#Emailid").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Please Enter a valid email address.</span>');
            hasError = true;
        }
 
       
        if(passw == '') {
            $("#Password").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Please enter your password.</span>');
            hasError = true;
          } 
          else if(!passreg.test(passw)) {
            $("#Password").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Enter a valid password.</span>');
            hasError = true;
        }
          if(hasError == true) { return false; }
    });


    $('#btn_Register').click(function() {  
  
        $(".error1").hide();
        var hasError = false;
        var emailpat = /^([\w-\.]+@([\w-]+\.)+[\w-]{2,4})?$/;
        var passpat=/^(?=.*\d)(?=.*[A-Z])(?=.*[a-z]).{8,}$/;
        var mobpat=/^([6-9]{1}[0-9]{9})$/;
      
        var name = $("#Name").val();
        var mobile = $("#Mobileno").val();
        var email = $("#Emailid").val();
        var location = $("#Location").val();
        var pass = $("#Password").val();
        if(name == '') {
            $("#Name").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Please Enter your First Name.</span>');
            hasError = true;
          } 
          if(location == '') {
            $("#Location").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Please Enter your First Name.</span>');
            hasError = true;
          } 
        if(email == '') {
            $("#Emailid").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Enter a your email address.</span>');
            hasError = true;
        }
 
        else if(!emailpat.test(email)) {
            $("#Emailid").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Enter a valid email address.</span>');
            hasError = true;
        }
 
       
        if(pass == '') {
            $("#Password").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Please Set Password.</span>');
            hasError = true;
          } 
          else if(!passpat.test(pass)) {
            $("#Password").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Password must be in(a-z,A-Z,min 8 digit,one special($,@)).</span>');
            hasError = true;
        }

        if(mobile == '') {
            $("#Mobileno").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Please Enter your Mobile no.</span>');
            hasError = true;
          } 
          else if(!mobpat.test(mobile)) {
            $("#Mobileno").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error1">*Mobile number must be 10 digit.</span>');
            hasError = true;
        }

       

          if(hasError == true) { return false; }
    });

    

});

