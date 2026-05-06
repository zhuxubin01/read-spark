import CoreData

final class CoreDataStack {
    static let shared = CoreDataStack()

    let container: NSPersistentContainer

    private init() {
        container = NSPersistentContainer(name: "ReadSpark")
        container.loadPersistentStores { _, error in
            if let error {
                assertionFailure("Core Data failed to load: \(error)")
            }
        }
        container.viewContext.automaticallyMergesChangesFromParent = true
    }
}
