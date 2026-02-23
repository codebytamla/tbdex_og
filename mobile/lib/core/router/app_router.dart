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
          pageBuilder: (context, state) => const NoTransitionPage(
            child: HomeScreen(),
          ),
        ),
        GoRoute(
          path: '/patients',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: PatientsScreen(),
          ),
        ),
        GoRoute(
          path: '/activity',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: ActivityScreen(),
          ),
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
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
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
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.home_outlined),
            selectedIcon: Icon(Icons.home_rounded),
            label: 'Home',
          ),
          NavigationDestination(
            icon: Icon(Icons.people_outline_rounded),
            selectedIcon: Icon(Icons.people_rounded),
            label: 'Patients',
          ),
          NavigationDestination(
            icon: Icon(Icons.assessment_outlined),
            selectedIcon: Icon(Icons.assessment_rounded),
            label: 'Activity',
          ),
        ],
      ),
    );
  }
}
