import mailchimp from '@mailchimp/mailchimp_marketing';

// Configure Mailchimp
mailchimp.setConfig({
  apiKey: import.meta.env.MAILCHIMP_API_KEY,
  server: import.meta.env.MAILCHIMP_SERVER_PREFIX, // e.g., 'us1', 'us2', etc.
});

export async function POST({ request }) {
  try {
    const { email } = await request.json();

    if (!email) {
      return new Response(
        JSON.stringify({ error: 'Email is required' }),
        { 
          status: 400,
          headers: { 'Content-Type': 'application/json' }
        }
      );
    }

    // Add subscriber to Mailchimp list
    const response = await mailchimp.lists.addListMember(
      import.meta.env.MAILCHIMP_AUDIENCE_ID,
      {
        email_address: email,
        status: 'subscribed',
        tags: ['early-access', 'loyaltywallet-beta'],
        merge_fields: {
          SOURCE: 'LoyaltyWallet Landing Page',
          SIGNUP_DATE: new Date().toISOString().split('T')[0]
        }
      }
    );

    return new Response(
      JSON.stringify({ 
        success: true, 
        message: 'Successfully subscribed to early access!',
        id: response.id 
      }),
      { 
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      }
    );

  } catch (error) {
    console.error('Mailchimp subscription error:', error);

    // Handle duplicate email addresses
    if (error.status === 400 && error.response?.body?.title === 'Member Exists') {
      return new Response(
        JSON.stringify({ 
          success: true, 
          message: 'You are already subscribed to our early access list!' 
        }),
        { 
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        }
      );
    }

    return new Response(
      JSON.stringify({ 
        error: 'Failed to subscribe. Please try again later.',
        details: error.message 
      }),
      { 
        status: 500,
        headers: { 'Content-Type': 'application/json' }
      }
    );
  }
}