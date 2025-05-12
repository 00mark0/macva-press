import { useState, useEffect } from 'react';
import { format, subDays, parseISO, isValid } from 'date-fns';
import fetch from './axios';

export default function DailyAnalytics() {
	const [analytics, setAnalytics] = useState([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [startDate, setStartDate] = useState(format(subDays(new Date(), 7), 'yyyy-MM-dd'));
	const [endDate, setEndDate] = useState(format(new Date(), 'yyyy-MM-dd'));
	const [limit, setLimit] = useState(30);

	/*function getCookie(name) {
		const value = `; ${document.cookie}`;
		const parts = value.split(`; ${name}=`);
		if (parts.length === 2) return parts.pop().split(';').shift();
		return null;
	}*/

	//const token = getCookie('access_token');

	const fetchAnalytics = async () => {
		setLoading(true);
		setError(null);

		try {
			const res = await fetch.get(`/api/admin/analytics?start_date=${startDate}&end_date=${endDate}&limit=${limit}`);

			// Handle null response by setting to empty array
			if (res.data === null) {
				setAnalytics([]);
			} else {
				setAnalytics(Array.isArray(res.data) ? res.data : []);
			}

		} catch (error) {
			console.error('Error fetching analytics:', error);
			setError('Failed to fetch analytics data');
			setAnalytics([]);
		} finally {
			setLoading(false);
		}
	};

	// Calculate totals for the summary card
	const calculateTotals = () => {
		if (!analytics || analytics.length === 0) return {
			views: 0, likes: 0, dislikes: 0, comments: 0, adsClicks: 0
		};

		return analytics.reduce((acc, day) => {
			return {
				views: acc.views + (day.total_views || 0),
				likes: acc.likes + (day.total_likes || 0),
				dislikes: acc.dislikes + (day.total_dislikes || 0),
				comments: acc.comments + (day.total_comments || 0),
				adsClicks: acc.adsClicks + (day.total_ads_clicks || 0),
			};
		}, { views: 0, likes: 0, dislikes: 0, comments: 0, adsClicks: 0 });
	};

	const handleSubmit = (e) => {
		e.preventDefault();
		fetchAnalytics();
	};

	const handleQuickDateRange = (days) => {
		const end = new Date();
		const start = subDays(end, days);

		setStartDate(format(start, 'yyyy-MM-dd'));
		setEndDate(format(end, 'yyyy-MM-dd'));

		// Trigger fetch immediately after state updates by using the values directly
		//setLoading(true);
		//setTimeout(() => {
		//fetchAnalytics();
		//}, 0);
	};

	// Format date for display
	const formatDate = (dateString) => {
		if (!dateString) return '';

		// Try to parse the date
		const date = parseISO(dateString);
		if (!isValid(date)) return dateString;

		return format(date, 'dd.MM.yyyy');
	};

	// Format numbers with thousands separators
	const formatNumber = (num) => {
		return num?.toLocaleString() || '0';
	};

	// Initialize on component mount
	useEffect(() => {
		fetchAnalytics();
	}, []);

	// Calculate totals for the summary
	const totals = calculateTotals();

	return (
		<div className="w-full min-h-screen dark:bg-black sm:p-8 p-4">
			<h1 className="text-black dark:text-white text-2xl font-bold mb-6">Dnevna analitika</h1>

			{/* Date range selector */}
			<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 mb-6">
				<form onSubmit={handleSubmit} className="flex flex-col flex-wrap space-y-4 sm:flex-row sm:space-y-0 sm:space-x-4 items-center">
					<div className="flex-1">
						<label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Početni datum
						</label>
						<input
							type="date"
							value={startDate}
							onChange={(e) => setStartDate(e.target.value)}
							className="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 bg-white dark:bg-gray-700 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<div className="flex-1">
						<label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Krajnji datum
						</label>
						<input
							type="date"
							value={endDate}
							onChange={(e) => setEndDate(e.target.value)}
							className="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 bg-white dark:bg-gray-700 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<div className="flex-1">
						<label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Limit rezultata
						</label>
						<select
							value={limit}
							onChange={(e) => setLimit(Number(e.target.value))}
							className="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 bg-white dark:bg-gray-700 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							<option value={7}>7 dana</option>
							<option value={30}>30 dana</option>
							<option value={90}>90 dana</option>
							<option value={180}>180 dana</option>
							<option value={365}>365 dana</option>
						</select>
					</div>

					<button
						type="submit"
						className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors sm:mt-6"
					>
						Prikaži
					</button>
				</form>

				{/* Quick date range buttons */}
				<div className="mt-4 flex flex-wrap gap-2">
					<button
						onClick={() => handleQuickDateRange(7)}
						className="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
					>
						Zadnjih 7 dana
					</button>
					<button
						onClick={() => handleQuickDateRange(30)}
						className="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
					>
						Zadnjih 30 dana
					</button>
					<button
						onClick={() => handleQuickDateRange(90)}
						className="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
					>
						Zadnjih 90 dana
					</button>
					<button
						onClick={() => handleQuickDateRange(180)}
						className="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
					>
						Zadnjih 180 dana
					</button>
				</div>
			</div>

			{/* Summary cards */}
			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
				<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
					<h3 className="text-gray-500 dark:text-gray-400 text-sm mb-1">Ukupno pregleda</h3>
					<p className="text-2xl font-bold text-black dark:text-white">{formatNumber(totals.views)}</p>
				</div>

				<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
					<h3 className="text-gray-500 dark:text-gray-400 text-sm mb-1">Ukupno lajkova</h3>
					<p className="text-2xl font-bold text-black dark:text-white">{formatNumber(totals.likes)}</p>
				</div>

				<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
					<h3 className="text-gray-500 dark:text-gray-400 text-sm mb-1">Ukupno dislajkova</h3>
					<p className="text-2xl font-bold text-black dark:text-white">{formatNumber(totals.dislikes)}</p>
				</div>

				<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
					<h3 className="text-gray-500 dark:text-gray-400 text-sm mb-1">Ukupno komentara</h3>
					<p className="text-2xl font-bold text-black dark:text-white">{formatNumber(totals.comments)}</p>
				</div>

				<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
					<h3 className="text-gray-500 dark:text-gray-400 text-sm mb-1">Ukupno klikova na oglase</h3>
					<p className="text-2xl font-bold text-black dark:text-white">{formatNumber(totals.adsClicks)}</p>
				</div>
			</div>

			{/* Loading state */}
			{loading ? (
				<div className="flex justify-center py-8">
					<div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			) : (
				<>
					{/* Error message */}
					{error && (
						<div className="bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 p-4 rounded-md mb-6">
							{error}. Molimo pokušajte ponovo kasnije.
						</div>
					)}

					{/* No data message */}
					{!error && analytics.length === 0 && (
						<div className="bg-gray-100 dark:bg-gray-800 p-6 rounded-lg text-center mb-6">
							<svg className="w-16 h-16 mx-auto text-gray-400 dark:text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
								<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
							</svg>
							<h3 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">
								Nema analitičkih podataka za izabrani period
							</h3>
							<p className="text-gray-600 dark:text-gray-400">
								Pokušajte sa drugim vremenskim periodom.
							</p>
						</div>
					)}

					{/* Analytics table */}
					{analytics.length > 0 && (
						<div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
							<div className="overflow-x-auto">
								<table className="w-full">
									<thead>
										<tr className="bg-gray-50 dark:bg-gray-700">
											<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Datum</th>
											<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Pregledi</th>
											<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Lajkovi</th>
											<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Dislajkovi</th>
											<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Komentari</th>
											<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Klikovi na oglase</th>
										</tr>
									</thead>
									<tbody className="divide-y divide-gray-200 dark:divide-gray-600">
										{analytics.map((day, index) => (
											<tr key={index} className="hover:bg-gray-50 dark:hover:bg-gray-700">
												<td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-100">
													{formatDate(day.analytics_date)}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
													{formatNumber(day.total_views)}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
													{formatNumber(day.total_likes)}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
													{formatNumber(day.total_dislikes)}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
													{formatNumber(day.total_comments)}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
													{formatNumber(day.total_ads_clicks)}
												</td>
											</tr>
										))}
									</tbody>
								</table>
							</div>
						</div>
					)}
				</>
			)}
		</div>
	);
}
