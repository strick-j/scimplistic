$(function() {
  
	// Get the form fields and hidden div
	var dbcheckbox = $('#EnableDatabase');
	var tlscheckbox = $('#EnableHTTPS')
	var hiddendb = $('#DatabaseGroup');
	var hiddentls = $('#TlsGroup');
	var authmethod = $('#AuthMethod');
	var hiddentoken = $('#OauthTokenField');
	var hiddencred = $('#ClientCredsGroup');
	
	// Initalize Form
	if (dbcheckbox.is(':checked')) {
		hiddendb.show();
	} else {
		hiddendb.hide();
	}

	if (authmethod.val() == "tokenauth") {
		// Show the hidden fields.
		hiddentoken.show();
		hiddencred.hide();
		$('#oathtoken').attr('required','');

		// Clear Client Creds if previously set
		$('#ClientID').removeAttr('required');
		$('#ClientSecret').removeAttr('required');
		$('#ClientAppID').removeAttr('required');
		$('#ClientID').val('');
		$('#ClientSecret').val('');
		$('#ClientAppID').val('');
	} else if ($(this).val() == "credauth") {
		hiddencred.show();
		hiddentoken.hide();
		// Set Client Info Required
		$('#ClientID').attr('required','');
		$('#ClientSecret').attr('required','');
		$('#ClientAppID').attr('required','');

		// Clear Token if previously set
		$('#oathtoken').removeAttr('required');
		$('#oathtoken').val('');
	} else {
		hiddentoken.hide();
		hiddencred.hide();
	}

	// Setup an event listener for when the state of the 
	// checkbox changes.
	dbcheckbox.change(function() {
		// Check to see if the checkbox is checked.
		// If it is, show the fields and populate the input.
		// If not, hide the fields.
		if (dbcheckbox.is(':checked')) {
			// Show the hidden fields.
			hiddendb.show();
			$('#DatabaseIP').attr('required','');
			$('#DatabasePort').attr('required','');
			$('#DatabaseName').attr('required','');
			$('#DatabaseUser').attr('required','');
			$('#DatabasePass').attr('required','');
		} else {
			// Make sure that the hidden fields are indeed
			// hidden.
			hiddendb.hide();
			$('#DatabaseIP').removeAttr('required');
			$('#DatabasePort').removeAttr('required');
			$('#DatabaseName').removeAttr('required');
			$('#DatabaseUser').removeAttr('required');
			$('#DatabasePass').removeAttr('required');

			// Clear fields incase user checks and unchecks
			$('#DatabaseIP').val('');
			$('#DatabasePort').val('');
			$('#DatabaseName').val('');
			$('#DatabaseUser').val('');
			$('#DatabasePass').val('');
		}
	});

	// Setup an event listener for when the state of the 
	// checkbox changes.
	tlscheckbox.change(function() {
		// Check to see if the checkbox is checked.
		// If it is, show the fields and populate the input.
		// If not, hide the fields.
		if (tlscheckbox.is(':checked')) {
			// Show the hidden fields.
			hiddentls.show();
			$('#ServerCertKey').attr('required','');
			$('#ServerCertFile').attr('required','');
		} else {
			// Make sure that the hidden fields are indeed
			// hidden.
			hiddentls.hide();

			// Clear fields incase user checks and unchecks
			$('#ServerCertKey').val('');
			$('#ServerCertFile').val('');
		}
	});

	// Setup an event listener for when the state of the 
	// checkbox changes for Auth Method.
	authmethod.change(function() {
		// Check to see if the checkbox is checked.
		// If it is, show the fields and populate the input.
		// If not, hide the fields.
		if ($(this).val() == "tokenauth") {
			// Show the hidden fields.
			hiddentoken.show();
			hiddencred.hide();
			$('#oathtoken').attr('required','');

			// Clear Client Creds if previously set
			$('#ClientID').removeAttr('required');
			$('#ClientSecret').removeAttr('required');
			$('#ClientAppID').removeAttr('required');
			$('#ClientID').val('');
			$('#ClientSecret').val('');
			$('#ClientAppID').val('');
		} else if ($(this).val() == "credauth") {
			hiddencred.show();
			hiddentoken.hide();
			// Set Client Info Required
			$('#ClientID').attr('required','');
			$('#ClientSecret').attr('required','');
			$('#ClientAppID').attr('required','');

			// Clear Token if previously set
			$('#oathtoken').removeAttr('required');
			$('#oathtoken').val('');
		} else {
			hiddentoken.hide();
			hiddencred.hide();

			// Clear Token if previously set
			$('#oathtoken').removeAttr('required');
			$('#oathtoken').val('');

			// Clear Client Creds if previously set
			$('#ClientID').removeAttr('required');
			$('#ClientSecret').removeAttr('required');
			$('#ClientAppID').removeAttr('required');
			$('#ClientID').val('');
			$('#ClientSecret').val('');
			$('#ClientAppID').val('');
		}
	});
	authmethod.trigger("change");
});