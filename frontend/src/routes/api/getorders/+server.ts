import * as db from '$lib/server/database';

export async function GET({ cookies }) {
	const sessionToken = cookies.get('session_token');

	const orders = await db.getOrders(sessionToken);

	if (!orders) {
		return new Response(JSON.stringify({ message: 'Error in getting orders' }), {
			status: 520
		});
	}

	return new Response(JSON.stringify({ orders }), {
		status: 200
	});
}
