import SwiftUI

@main
struct ReadSparkApp: App {
    let persistenceController = CoreDataStack.shared

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(\.managedObjectContext, persistenceController.container.viewContext)
        }
    }
}
