import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/app_theme.dart';
import '../../core/data/mock_data.dart';
import '../../core/models/models.dart';
import '../../../main.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  String _greeting() {
    final hour = DateTime.now().hour;
    if (hour < 12) return 'Good morning';
    if (hour < 17) return 'Good afternoon';
    return 'Good evening';
  }

  String _relativeTime(DateTime timestamp) {
    final diff = DateTime.now().difference(timestamp);
    if (diff.inMinutes < 1) return 'just now';
    if (diff.inMinutes < 60) return '${diff.inMinutes} min ago';
    if (diff.inHours < 24) return '${diff.inHours}h ago';
    return '${diff.inDays}d ago';
  }

  IconData _syncIcon({
    required SyncStatus syncStatus,
    required bool isLastSyncStale,
  }) {
    if (!syncStatus.isOnline) return Icons.cloud_off_rounded;
    if (syncStatus.failedItemCount > 0) return Icons.sync_problem_rounded;
    if (syncStatus.pendingTotal > 0 || isLastSyncStale) {
      return Icons.sync_rounded;
    }
    return Icons.check_circle_rounded;
  }

  String _syncTitle({
    required SyncStatus syncStatus,
    required bool isLastSyncStale,
  }) {
    if (!syncStatus.isOnline) return 'Offline: sync paused';
    if (syncStatus.failedItemCount > 0) return 'Sync needs attention';
    if (syncStatus.pendingUploadCount > 0) return 'Tap to sync data';
    if (syncStatus.pendingDownloadCount > 0) return 'Updates ready to sync';
    if (isLastSyncStale) return 'Sync recommended';
    return 'Synced';
  }

  String _syncSubtitle({
    required SyncStatus syncStatus,
    required bool isLastSyncStale,
  }) {
    if (!syncStatus.isOnline) return 'Connect to the internet to sync';
    if (syncStatus.failedItemCount > 0) {
      final suffix = syncStatus.failedItemCount == 1 ? '' : 's';
      return '${syncStatus.failedItemCount} failed item$suffix. Tap to review';
    }
    if (syncStatus.pendingTotal > 0) {
      final segments = <String>[];
      if (syncStatus.pendingUploadCount > 0) {
        final suffix = syncStatus.pendingUploadCount == 1 ? '' : 's';
        segments.add('${syncStatus.pendingUploadCount} upload$suffix');
      }
      if (syncStatus.pendingDownloadCount > 0) {
        final suffix = syncStatus.pendingDownloadCount == 1 ? '' : 's';
        segments.add('${syncStatus.pendingDownloadCount} download$suffix');
      }
      return '${segments.join(' • ')} pending';
    }
    if (isLastSyncStale) {
      return 'Last synced ${_relativeTime(syncStatus.lastSyncedAt)}';
    }
    return 'Last synced ${_relativeTime(syncStatus.lastSyncedAt)}';
  }

  @override
  Widget build(BuildContext context) {
    final tasks = MockData.todaysTasks;
    const previewTaskCount = 3;
    final previewTasks = tasks.take(previewTaskCount).toList();
    final hasMoreTasks = tasks.length > previewTaskCount;
    final assessmentCount = tasks
        .where((t) => t.status == 'pending' || t.status == 'in_progress')
        .length;
    final referralCount = MockData.assessmentHistory
        .where((a) => a.referralCreated)
        .length;
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final syncStatus = MockData.syncStatus;
    final isLastSyncStale =
        DateTime.now().difference(syncStatus.lastSyncedAt) >
        const Duration(hours: 8);
    final showSyncCard =
        !syncStatus.isOnline ||
        syncStatus.pendingTotal > 0 ||
        syncStatus.failedItemCount > 0 ||
        isLastSyncStale;
    final syncTitle = _syncTitle(
      syncStatus: syncStatus,
      isLastSyncStale: isLastSyncStale,
    );
    final syncSubtitle = _syncSubtitle(
      syncStatus: syncStatus,
      isLastSyncStale: isLastSyncStale,
    );
    final syncIcon = _syncIcon(
      syncStatus: syncStatus,
      isLastSyncStale: isLastSyncStale,
    );

    return Scaffold(
      backgroundColor: isDark ? AppColorsDark.background : AppColors.background,
      body: SafeArea(
        bottom: false,
        child: CustomScrollView(
          slivers: [
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(22, 16, 22, 0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // ── Header ──────────────────────────────────────────────
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                _greeting(),
                                style: Theme.of(context).textTheme.bodyMedium
                                    ?.copyWith(fontWeight: FontWeight.w500),
                              ),
                              const SizedBox(height: 2),
                              Text(
                                'Overview',
                                style: Theme.of(context).textTheme.displaySmall
                                    ?.copyWith(
                                      fontWeight: FontWeight.w800,
                                      letterSpacing: -0.5,
                                      height: 1.1,
                                    ),
                              ),
                            ],
                          ),
                        ),
                        ValueListenableBuilder<ThemeMode>(
                          valueListenable: themeNotifier,
                          builder: (context, mode, child) => IconButton(
                            onPressed: () {
                              themeNotifier.value = mode == ThemeMode.light
                                  ? ThemeMode.dark
                                  : ThemeMode.light;
                            },
                            icon: Icon(
                              mode == ThemeMode.light
                                  ? Icons.dark_mode
                                  : Icons.light_mode,
                              size: 22,
                              color: isDark
                                  ? AppColorsDark.textSecondary
                                  : AppColors.textSecondary,
                            ),
                          ),
                        ),
                        const SizedBox(width: 2),
                        // User avatar
                        Container(
                          width: 40,
                          height: 40,
                          decoration: BoxDecoration(
                            color: isDark
                                ? AppColorsDark.surfaceDim
                                : const Color(0xFFE9EEF3),
                            borderRadius: BorderRadius.circular(999),
                          ),
                          alignment: Alignment.center,
                          child: Text(
                            'TD',
                            style: TextStyle(
                              color: isDark
                                  ? AppColorsDark.textSecondary
                                  : AppColors.textSecondary,
                              fontWeight: FontWeight.w700,
                              fontSize: 15,
                            ),
                          ),
                        ),
                      ],
                    ),

                    const SizedBox(height: 18),

                    // ── Sync Banner ─────────────────────────────────────────
                    if (showSyncCard)
                      Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: () => context.push('/sync'),
                          borderRadius: BorderRadius.circular(14),
                          child: Ink(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 14,
                              vertical: 7,
                            ),
                            decoration: BoxDecoration(
                              color: isDark
                                  ? AppColorsDark.surfaceDim
                                  : AppColors.surfaceDim,
                              borderRadius: BorderRadius.circular(14),
                            ),
                            child: Row(
                              children: [
                                Container(
                                  width: 28,
                                  height: 28,
                                  decoration: BoxDecoration(
                                    color: isDark
                                        ? AppColorsDark.surface
                                        : const Color(0xFFE9EEF3),
                                    borderRadius: BorderRadius.circular(10),
                                  ),
                                  child: Icon(
                                    syncIcon,
                                    color: isDark
                                        ? AppColorsDark.textSecondary
                                        : const Color(0xFF556275),
                                    size: 18,
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      Text(
                                        syncTitle,
                                        style: TextStyle(
                                          color: isDark
                                              ? AppColorsDark.textPrimary
                                              : AppColors.textPrimary,
                                          fontWeight: FontWeight.w600,
                                          fontSize: 15,
                                        ),
                                      ),
                                      const SizedBox(height: 1),
                                      Text(
                                        syncSubtitle,
                                        style: TextStyle(
                                          color: isDark
                                              ? AppColorsDark.textSecondary
                                              : AppColors.textSecondary,
                                          fontWeight: FontWeight.w400,
                                          fontSize: 13,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                Icon(
                                  Icons.chevron_right_rounded,
                                  color: isDark
                                      ? AppColorsDark.textSecondary
                                      : AppColors.textSecondary,
                                  size: 20,
                                ),
                              ],
                            ),
                          ),
                        ),
                      )
                    else
                      Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: () => context.push('/sync'),
                          borderRadius: BorderRadius.circular(10),
                          child: Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 2,
                              vertical: 3,
                            ),
                            child: Row(
                              children: [
                                Icon(
                                  Icons.check_circle_rounded,
                                  size: 18,
                                  color: isDark
                                      ? const Color(0xFF95C0B0)
                                      : const Color(0xFF3F7A5E),
                                ),
                                const SizedBox(width: 8),
                                Expanded(
                                  child: Text(
                                    'Synced ${_relativeTime(syncStatus.lastSyncedAt)}',
                                    style: TextStyle(
                                      color: isDark
                                          ? AppColorsDark.textSecondary
                                          : AppColors.textSecondary,
                                      fontWeight: FontWeight.w500,
                                      fontSize: 14,
                                    ),
                                  ),
                                ),
                                Text(
                                  'Sync center',
                                  style: TextStyle(
                                    color: isDark
                                        ? AppColorsDark.textTertiary
                                        : AppColors.textTertiary,
                                    fontWeight: FontWeight.w500,
                                    fontSize: 13,
                                  ),
                                ),
                                Icon(
                                  Icons.chevron_right_rounded,
                                  size: 18,
                                  color: isDark
                                      ? AppColorsDark.textTertiary
                                      : AppColors.textTertiary,
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),

                    const SizedBox(height: 18),

                    // ── Stats Row ────────────────────────────────────────────
                    Row(
                      children: [
                        Expanded(
                          child: _MiniStatCard(
                            label: 'ASSESSMENTS',
                            value: '$assessmentCount',
                            icon: Icons.description,
                            iconColor: const Color(0xFF5E6B7D),
                            iconColorDark: const Color(0xFFA5B1C0),
                            iconBg: const Color(0xFFF1F4F8),
                            iconBgDark: const Color(
                              0xFF2B3441,
                            ).withValues(alpha: 0.45),
                            isDark: isDark,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _MiniStatCard(
                            label: 'REFERRALS',
                            value: '$referralCount',
                            icon: Icons.person_add,
                            iconColor: const Color(0xFF6D6256),
                            iconColorDark: const Color(0xFFC2B5A5),
                            iconBg: const Color(0xFFF5F1EA),
                            iconBgDark: const Color(
                              0xFF3A3228,
                            ).withValues(alpha: 0.45),
                            isDark: isDark,
                          ),
                        ),
                      ],
                    ),

                    const SizedBox(height: 36),

                    // ── Recent Assessments header ────────────────────────────
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Text(
                          'Recents',
                          style: Theme.of(context).textTheme.headlineSmall
                              ?.copyWith(fontWeight: FontWeight.w700),
                        ),
                        if (hasMoreTasks)
                          InkWell(
                            onTap: () => context.go('/activity'),
                            borderRadius: BorderRadius.circular(8),
                            child: Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 4,
                                vertical: 2,
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Text(
                                    'See All',
                                    style: TextStyle(
                                      color: Theme.of(
                                        context,
                                      ).colorScheme.primary,
                                      fontWeight: FontWeight.w600,
                                      fontSize: 16,
                                    ),
                                  ),
                                  Icon(
                                    Icons.chevron_right_rounded,
                                    size: 18,
                                    color: Theme.of(
                                      context,
                                    ).colorScheme.primary,
                                  ),
                                ],
                              ),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 12),
                  ],
                ),
              ),
            ),

            // ── Task list ────────────────────────────────────────────────────
            SliverPadding(
              padding: const EdgeInsets.symmetric(horizontal: 22),
              sliver: SliverList.separated(
                itemCount: previewTasks.length,
                separatorBuilder: (_, _) => const SizedBox(height: 10),
                itemBuilder: (context, index) {
                  final task = previewTasks[index];
                  return _TaskRow(
                    task: task,
                    onTap: () => context.push('/patients/${task.patientId}'),
                  );
                },
              ),
            ),

            const SliverToBoxAdapter(child: SizedBox(height: 100)),
          ],
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Mini Stat Card (Assessments / Referrals)
// ─────────────────────────────────────────────────────────────────────────────
class _MiniStatCard extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color iconColor;
  final Color iconColorDark;
  final Color iconBg;
  final Color iconBgDark;
  final bool isDark;

  const _MiniStatCard({
    required this.label,
    required this.value,
    required this.icon,
    required this.iconColor,
    required this.iconColorDark,
    required this.iconBg,
    required this.iconBgDark,
    required this.isDark,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveColor = isDark ? iconColorDark : iconColor;
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? AppColorsDark.surface : AppColors.surface,
        borderRadius: BorderRadius.circular(18),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: isDark ? 0.08 : 0.025),
            blurRadius: 8,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: isDark ? iconBgDark : iconBg,
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: effectiveColor, size: 18),
          ),
          const SizedBox(height: 16),
          Text(
            value,
            style: TextStyle(
              fontSize: 34,
              fontWeight: FontWeight.w700,
              color: isDark ? AppColorsDark.textPrimary : AppColors.textPrimary,
              height: 1,
            ),
          ),
          const SizedBox(height: 5),
          Text(
            label,
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w600,
              color: isDark
                  ? AppColorsDark.textSecondary
                  : AppColors.textSecondary,
              letterSpacing: 0.9,
            ),
          ),
        ],
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Task Row
// ─────────────────────────────────────────────────────────────────────────────
class _TaskRow extends StatelessWidget {
  final TaskItem task;
  final VoidCallback onTap;

  const _TaskRow({required this.task, required this.onTap});

  String _initials(String name) {
    final parts = name.split(' ');
    if (parts.length >= 2) return '${parts[0][0]}${parts[1][0]}'.toUpperCase();
    return name[0].toUpperCase();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        borderRadius: BorderRadius.circular(20),
        onTap: onTap,
        child: Ink(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          decoration: BoxDecoration(
            color: isDark ? AppColorsDark.surface : AppColors.surface,
            borderRadius: BorderRadius.circular(20),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: isDark ? 0.08 : 0.025),
                blurRadius: 8,
                offset: const Offset(0, 1),
              ),
            ],
          ),
          child: Row(
            children: [
              Container(
                width: 36,
                height: 36,
                decoration: BoxDecoration(
                  color: isDark
                      ? AppColorsDark.surfaceDim
                      : AppColors.surfaceDim,
                  borderRadius: BorderRadius.circular(10),
                ),
                alignment: Alignment.center,
                child: Text(
                  _initials(task.patientName),
                  style: TextStyle(
                    color: isDark
                        ? AppColorsDark.textSecondary
                        : AppColors.textSecondary,
                    fontWeight: FontWeight.w600,
                    fontSize: 13,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      task.patientName,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 3),
                    Text(
                      '${task.village} • ${task.time}',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.w500,
                        color: isDark
                            ? AppColorsDark.textSecondary
                            : AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 10),
              _StatusChip(status: task.status),
              const SizedBox(width: 6),
              Icon(
                Icons.chevron_right_rounded,
                size: 18,
                color: isDark
                    ? AppColorsDark.textTertiary
                    : AppColors.textTertiary,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Status Chip
// ─────────────────────────────────────────────────────────────────────────────
class _StatusChip extends StatelessWidget {
  final String status;
  const _StatusChip({required this.status});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final style = _StatusStyle.resolve(status, isDark);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: style.background,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        style.label,
        style: TextStyle(
          color: style.foreground,
          fontSize: 13,
          fontWeight: FontWeight.w600,
          letterSpacing: 0.3,
        ),
      ),
    );
  }
}

class _StatusStyle {
  final Color foreground;
  final Color background;
  final String label;

  const _StatusStyle({
    required this.foreground,
    required this.background,
    required this.label,
  });

  static _StatusStyle resolve(String status, bool isDark) {
    switch (status) {
      case 'pending':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFFE6D3A0)
              : const Color(0xFF7A5A2A),
          background: isDark
              ? AppColorsDark.warningLight.withValues(alpha: 0.3)
              : AppColors.warningLight.withValues(alpha: 0.7),
          label: 'PENDING',
        );
      case 'in_progress':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFFB6C7E5)
              : const Color(0xFF3E5A7A),
          background: isDark
              ? const Color(0xFF273246).withValues(alpha: 0.55)
              : const Color(0xFFEAF0FA),
          label: 'IN PROGRESS',
        );
      case 'done':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFFA7F3D0)
              : const Color(0xFF115E59),
          background: isDark
              ? AppColorsDark.successLight.withValues(alpha: 0.32)
              : AppColors.successLight.withValues(alpha: 0.75),
          label: 'DONE',
        );
      default:
        return _StatusStyle(
          foreground: isDark
              ? AppColorsDark.textSecondary
              : AppColors.textSecondary,
          background: isDark ? AppColorsDark.surfaceDim : AppColors.surfaceDim,
          label: status.toUpperCase(),
        );
    }
  }
}
