import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/app_theme.dart';
import '../../core/data/mock_data.dart';

class SyncCenterScreen extends StatelessWidget {
  const SyncCenterScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Sync Center'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_rounded),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Connection status
            _ConnectionBanner(isConnected: false),

            const SizedBox(height: 20),

            // Summary cards
            Row(
              children: [
                Expanded(
                  child: _SyncSummaryCard(
                    icon: Icons.cloud_upload_rounded,
                    iconColor: AppColors.primary,
                    title: 'Pending Uploads',
                    detail: '8 assessments\n3 referrals\n2 new patients',
                    badge: '13',
                  ),
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: _SyncSummaryCard(
                    icon: Icons.check_circle_rounded,
                    iconColor: AppColors.success,
                    title: 'Last Sync',
                    detail: 'Feb 23, 2026\n8:30 AM',
                  ),
                ),
              ],
            ),

            const SizedBox(height: 10),

            _SyncSummaryCard(
              icon: Icons.download_rounded,
              iconColor: AppColors.info,
              title: 'Download Queue',
              detail: '2 template updates available',
              isFullWidth: true,
            ),

            const SizedBox(height: 24),

            // Sync history
            Text(
              'Recent Sync History',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 12),

            ...MockData.syncHistory.map(
              (event) => _SyncHistoryItem(event: event),
            ),

            const SizedBox(height: 32),

            // Sync button
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: null, // Disabled when offline
                icon: const Icon(Icons.sync_rounded),
                label: const Text('Sync Now'),
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  disabledBackgroundColor: AppColors.surfaceDim,
                  disabledForegroundColor: AppColors.textTertiary,
                ),
              ),
            ),

            const SizedBox(height: 8),

            Center(
              child: Text(
                'Connect to the internet to sync',
                style: TextStyle(
                  fontSize: 12,
                  color: AppColors.textTertiary,
                ),
              ),
            ),

            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }
}

class _ConnectionBanner extends StatelessWidget {
  final bool isConnected;

  const _ConnectionBanner({required this.isConnected});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isConnected ? AppColors.successLight : AppColors.surfaceDim,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isConnected
              ? AppColors.success.withValues(alpha: 0.3)
              : AppColors.border,
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 12,
            height: 12,
            decoration: BoxDecoration(
              color: isConnected ? AppColors.success : AppColors.textTertiary,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 12),
          Text(
            isConnected
                ? 'Connected via Mobile Data'
                : 'No Internet Connection',
            style: TextStyle(
              fontSize: 15,
              fontWeight: FontWeight.w600,
              color: isConnected
                  ? AppColors.success
                  : AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}

class _SyncSummaryCard extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final String title;
  final String detail;
  final String? badge;
  final bool isFullWidth;

  const _SyncSummaryCard({
    required this.icon,
    required this.iconColor,
    required this.title,
    required this.detail,
    this.badge,
    this.isFullWidth = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: isFullWidth ? double.infinity : null,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(icon, size: 20, color: iconColor),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  title,
                  style: TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textSecondary,
                  ),
                ),
              ),
              if (badge != null)
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(
                    color: AppColors.primary,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    badge!,
                    style: const TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.w700,
                      color: Colors.white,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 10),
          Text(
            detail,
            style: TextStyle(
              fontSize: 13,
              color: AppColors.textPrimary,
              height: 1.5,
            ),
          ),
        ],
      ),
    );
  }
}

class _SyncHistoryItem extends StatelessWidget {
  final dynamic event;

  const _SyncHistoryItem({required this.event});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(
            event.success
                ? Icons.check_circle_rounded
                : Icons.error_rounded,
            size: 18,
            color: event.success ? AppColors.success : AppColors.danger,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  event.message,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: event.success
                        ? AppColors.textPrimary
                        : AppColors.danger,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  event.time,
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ],
            ),
          ),
          if (event.retryAvailable)
            TextButton(
              onPressed: () {},
              style: TextButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 8),
                minimumSize: Size.zero,
                tapTargetSize: MaterialTapTargetSize.shrinkWrap,
              ),
              child: const Text(
                'Retry',
                style: TextStyle(fontSize: 13),
              ),
            ),
        ],
      ),
    );
  }
}
