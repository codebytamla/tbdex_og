import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../features/home/home_screen.dart';
import '../../features/patients/patients_screen.dart';
import '../../features/patients/patient_detail_screen.dart';
import '../../features/activity/activity_screen.dart';
import '../../features/assessment/assessment_wizard_screen.dart';
import '../../features/assessment/assessment_result_screen.dart';
import '../../features/sync/sync_center_screen.dart';
import '../models/models.dart';
import '../data/mock_data.dart';
import '../theme/app_theme.dart';
import '../widgets/liquid_glass_nav_bar.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

final GoRouter appRouter = GoRouter(
  navigatorKey: _rootNavigatorKey,
  initialLocation: '/home',
  routes: [
    ShellRoute(
      navigatorKey: _shellNavigatorKey,
      builder: (context, state, child) => AppShell(child: child),
      routes: [
        GoRoute(
          path: '/home',
          pageBuilder: (context, state) =>
              const NoTransitionPage(child: HomeScreen()),
        ),
        GoRoute(
          path: '/patients',
          pageBuilder: (context, state) =>
              const NoTransitionPage(child: PatientsScreen()),
        ),
        GoRoute(
          path: '/activity',
          pageBuilder: (context, state) =>
              const NoTransitionPage(child: ActivityScreen()),
        ),
      ],
    ),
    GoRoute(
      path: '/patients/:id',
      parentNavigatorKey: _rootNavigatorKey,
      builder: (context, state) {
        final patientId = state.pathParameters['id']!;
        return PatientDetailScreen(patientId: patientId);
      },
    ),
    GoRoute(
      path: '/assessment/:patientId',
      parentNavigatorKey: _rootNavigatorKey,
      builder: (context, state) {
        final patientId = state.pathParameters['patientId']!;
        return AssessmentWizardScreen(patientId: patientId);
      },
    ),
    GoRoute(
      path: '/assessment-result',
      parentNavigatorKey: _rootNavigatorKey,
      builder: (context, state) {
        final data = state.extra as Map<String, dynamic>;
        return AssessmentResultScreen(
          assessmentData: data['assessmentData'] as AssessmentData,
          patientId: data['patientId'] as String,
        );
      },
    ),
    GoRoute(
      path: '/sync',
      parentNavigatorKey: _rootNavigatorKey,
      builder: (context, state) => const SyncCenterScreen(),
    ),
  ],
);

class AppShell extends StatelessWidget {
  final Widget child;

  const AppShell({super.key, required this.child});

  int _currentIndex(BuildContext context) {
    final location = GoRouterState.of(context).uri.path;
    if (location.startsWith('/patients')) return 1;
    if (location.startsWith('/activity')) return 2;
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      extendBody: true,
      body: Stack(
        children: [
          child,
          LiquidGlassNavBar(
            selectedIndex: _currentIndex(context),
            onDestinationSelected: (index) {
              switch (index) {
                case 0:
                  context.go('/home');
                case 1:
                  context.go('/patients');
                case 2:
                  context.go('/activity');
              }
            },
            items: const [
              LiquidGlassNavItem(
                icon: Icons.home_outlined,
                selectedIcon: Icons.home_rounded,
                label: 'Home',
              ),
              LiquidGlassNavItem(
                icon: Icons.people_outline_rounded,
                selectedIcon: Icons.people_rounded,
                label: 'Patients',
              ),
              LiquidGlassNavItem(
                icon: Icons.assessment_outlined,
                selectedIcon: Icons.assessment_rounded,
                label: 'Activity',
              ),
            ],
          ),
          // Floating FAB circle beside the pill nav bar
          Positioned(
            right: 20,
            bottom: 20,
            child: GestureDetector(
              onTap: () => _showPatientPicker(context),
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                width: 62,
                height: 62,
                decoration: BoxDecoration(
                  color: isDark ? AppColorsDark.primary : AppColors.primary,
                  borderRadius: BorderRadius.circular(999),
                  boxShadow: [
                    BoxShadow(
                      color:
                          (isDark ? AppColorsDark.primary : AppColors.primary)
                              .withValues(alpha: 0.35),
                      blurRadius: 18,
                      offset: const Offset(0, 7),
                    ),
                  ],
                ),
                child: const Icon(
                  Icons.add_rounded,
                  color: Colors.white,
                  size: 28,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _showPatientPicker(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => Container(
        height: MediaQuery.of(ctx).size.height * 0.7,
        decoration: BoxDecoration(
          color: isDark ? AppColorsDark.surface : AppColors.surface,
          borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
        ),
        child: Column(
          children: [
            const SizedBox(height: 8),
            Container(
              width: 36,
              height: 4,
              decoration: BoxDecoration(
                color: isDark ? AppColorsDark.border : AppColors.border,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Select Patient',
                    style: Theme.of(ctx).textTheme.headlineSmall,
                  ),
                  const SizedBox(height: 12),
                  const TextField(
                    decoration: InputDecoration(
                      hintText: 'Search patients...',
                      prefixIcon: Icon(Icons.search_rounded),
                    ),
                  ),
                ],
              ),
            ),
            const Divider(height: 1),
            Expanded(
              child: ListView.builder(
                padding: const EdgeInsets.symmetric(vertical: 8),
                itemCount: MockData.patients.length,
                itemBuilder: (ctx, index) {
                  final patient = MockData.patients[index];
                  return ListTile(
                    leading: CircleAvatar(
                      backgroundColor: isDark
                          ? AppColorsDark.primaryLight.withValues(alpha: 0.35)
                          : AppColors.primaryLight,
                      child: Text(
                        patient.initials,
                        style: TextStyle(
                          color: ctx.appPrimary,
                          fontWeight: FontWeight.w600,
                          fontSize: 13,
                        ),
                      ),
                    ),
                    title: Text(
                      patient.fullName,
                      style: const TextStyle(fontWeight: FontWeight.w500),
                    ),
                    subtitle: Text(
                      '${patient.village ?? ''} · ${patient.mrn ?? ''}',
                    ),
                    trailing: Icon(
                      Icons.arrow_forward_ios_rounded,
                      size: 14,
                      color: ctx.appTextTertiary,
                    ),
                    onTap: () {
                      Navigator.pop(ctx);
                      context.push('/assessment/${patient.id}');
                    },
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
